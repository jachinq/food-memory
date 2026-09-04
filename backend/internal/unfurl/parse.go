package unfurl

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type pageMeta struct {
	OGTitle       string
	OGImage       string
	OGSiteName    string
	TwitterTitle  string
	TwitterImage  string
	HTMLTitle     string
	RecipeName    string
	RecipeImage   string
	Ingredients   string
	CookMinutes   *int
}

var siteTitleSuffixes = []string{
	" - 小红书",
	" - YouTube",
	" - TikTok",
	" | 下厨房",
	" - 下厨房",
	" | Bilibili",
	" - 哔哩哔哩",
}

func parseHTML(r io.Reader) pageMeta {
	limited := io.LimitReader(r, maxBodyBytes)
	doc, err := html.Parse(limited)
	if err != nil {
		return pageMeta{}
	}
	meta := pageMeta{}
	var jsonLDs []string
	var inHead bool
	var titleBuf strings.Builder
	var inTitle bool

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch strings.ToLower(n.Data) {
			case "head":
				inHead = true
			case "body":
				if jsonLDs == nil {
					// continue; some sites put ld+json in body
				}
			case "title":
				if inHead && meta.HTMLTitle == "" {
					inTitle = true
					titleBuf.Reset()
				}
			case "meta":
				collectMeta(n, &meta)
			case "script":
				if attr(n, "type") == "application/ld+json" {
					jsonLDs = append(jsonLDs, strings.TrimSpace(textContent(n)))
				}
			}
		}
		if n.Type == html.TextNode && inTitle {
			titleBuf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
		if n.Type == html.ElementNode && strings.ToLower(n.Data) == "title" && inTitle {
			meta.HTMLTitle = strings.TrimSpace(titleBuf.String())
			inTitle = false
		}
	}
	walk(doc)

	for _, raw := range jsonLDs {
		applyRecipeJSONLD(raw, &meta)
	}
	return meta
}

func collectMeta(n *html.Node, meta *pageMeta) {
	prop := strings.ToLower(attr(n, "property"))
	name := strings.ToLower(attr(n, "name"))
	content := strings.TrimSpace(attr(n, "content"))
	if content == "" {
		return
	}
	key := prop
	if key == "" {
		key = name
	}
	switch key {
	case "og:title":
		if meta.OGTitle == "" {
			meta.OGTitle = content
		}
	case "og:image", "og:image:url", "og:image:secure_url":
		if meta.OGImage == "" {
			meta.OGImage = content
		}
	case "og:site_name":
		if meta.OGSiteName == "" {
			meta.OGSiteName = content
		}
	case "twitter:title":
		if meta.TwitterTitle == "" {
			meta.TwitterTitle = content
		}
	case "twitter:image", "twitter:image:src":
		if meta.TwitterImage == "" {
			meta.TwitterImage = content
		}
	}
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}

func textContent(n *html.Node) string {
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return b.String()
}

func cleanTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if i := strings.Index(s, " | "); i > 0 {
		s = strings.TrimSpace(s[:i])
	}
	for _, suf := range siteTitleSuffixes {
		s = strings.TrimSuffix(s, suf)
	}
	return strings.TrimSpace(s)
}

func applyRecipeJSONLD(raw string, meta *pageMeta) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}
	var v any
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		// some pages wrap with HTML comments or trailing commas — try a looser cut
		raw = strings.TrimPrefix(raw, "<!--")
		raw = strings.TrimSuffix(raw, "-->")
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return
		}
	}
	findRecipe(v, meta)
}

func findRecipe(v any, meta *pageMeta) {
	switch t := v.(type) {
	case []any:
		for _, item := range t {
			findRecipe(item, meta)
		}
	case map[string]any:
		if hasType(t["@type"], "Recipe") {
			applyRecipeMap(t, meta)
		}
		if g, ok := t["@graph"]; ok {
			findRecipe(g, meta)
		}
	}
}

func hasType(v any, want string) bool {
	switch t := v.(type) {
	case string:
		return typeName(t) == want
	case []any:
		for _, item := range t {
			if s, ok := item.(string); ok && typeName(s) == want {
				return true
			}
		}
	}
	return false
}

func typeName(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "/"); i >= 0 {
		s = s[i+1:]
	}
	return s
}

func applyRecipeMap(m map[string]any, meta *pageMeta) {
	if meta.RecipeName == "" {
		meta.RecipeName = strings.TrimSpace(asString(m["name"]))
	}
	if meta.RecipeImage == "" {
		meta.RecipeImage = firstImage(m["image"])
	}
	if meta.Ingredients == "" {
		meta.Ingredients = joinIngredients(m["recipeIngredient"])
	}
	if meta.CookMinutes == nil {
		if mins := parseDurationMinutes(asString(m["cookTime"])); mins != nil {
			meta.CookMinutes = mins
		} else if mins := parseDurationMinutes(asString(m["totalTime"])); mins != nil {
			meta.CookMinutes = mins
		}
	}
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case map[string]any:
		if s, ok := t["@value"].(string); ok {
			return s
		}
		if s, ok := t["name"].(string); ok {
			return s
		}
	}
	return ""
}

func firstImage(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case []any:
		for _, item := range t {
			if s := firstImage(item); s != "" {
				return s
			}
		}
	case map[string]any:
		if s := strings.TrimSpace(asString(t["url"])); s != "" {
			return s
		}
		if s := strings.TrimSpace(asString(t["contentUrl"])); s != "" {
			return s
		}
	}
	return ""
}

func joinIngredients(v any) string {
	var parts []string
	switch t := v.(type) {
	case string:
		parts = []string{strings.TrimSpace(t)}
	case []any:
		for _, item := range t {
			s := strings.TrimSpace(asString(item))
			if s != "" {
				parts = append(parts, s)
			}
		}
	}
	if len(parts) > 12 {
		parts = parts[:12]
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, "，")
}

var isoDuration = regexp.MustCompile(`(?i)^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

func parseDurationMinutes(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	m := isoDuration.FindStringSubmatch(s)
	if m == nil {
		return nil
	}
	days := atoiZero(m[1])
	hours := atoiZero(m[2])
	mins := atoiZero(m[3])
	secs := 0
	if m[4] != "" {
		f, err := strconv.ParseFloat(m[4], 64)
		if err == nil {
			secs = int(f + 0.5)
		}
	}
	total := days*24*60 + hours*60 + mins + (secs+59)/60
	if days == 0 && hours == 0 && mins == 0 && secs == 0 {
		return nil
	}
	if secs > 0 && mins == 0 && hours == 0 && days == 0 && total == 0 {
		total = 1
	}
	return &total
}

func atoiZero(s string) int {
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

func pickName(meta pageMeta, oembedTitle string) string {
	for _, s := range []string{meta.RecipeName, oembedTitle, meta.OGTitle, meta.TwitterTitle, meta.HTMLTitle} {
		if v := cleanTitle(s); v != "" {
			return v
		}
	}
	return ""
}

func pickCover(meta pageMeta, oembedThumb string) string {
	for _, s := range []string{meta.RecipeImage, oembedThumb, meta.OGImage, meta.TwitterImage} {
		if strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func pickPlatform(hostPlatform, ogSite, oembedProvider string) string {
	if hostPlatform != "" {
		return hostPlatform
	}
	if s := strings.TrimSpace(ogSite); s != "" {
		return s
	}
	return strings.TrimSpace(oembedProvider)
}

func readLimited(r io.Reader) ([]byte, error) {
	return io.ReadAll(io.LimitReader(r, maxBodyBytes))
}

func looksLikeHTML(b []byte) bool {
	s := bytes.TrimSpace(b)
	if len(s) == 0 {
		return false
	}
	if s[0] == '<' {
		return true
	}
	return bytes.Contains(bytes.ToLower(s[:min(len(s), 256)]), []byte("<html"))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
