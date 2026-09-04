package unfurl

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlatformFromHost(t *testing.T) {
	cases := map[string]string{
		"www.xiaohongshu.com":     "小红书",
		"xhslink.com":             "小红书",
		"www.douyin.com":          "抖音",
		"b23.tv":                  "B站",
		"www.bilibili.com:443":    "B站",
		"weibo.cn":                "微博",
		"mp.weixin.qq.com":        "公众号",
		"youtu.be":                "YouTube",
		"www.youtube.com":         "YouTube",
		"www.tiktok.com":          "TikTok",
		"www.xiachufang.com":      "下厨房",
		"example.com":             "",
	}
	for host, want := range cases {
		if got := PlatformFromHost(host); got != want {
			t.Errorf("PlatformFromHost(%q)=%q want %q", host, got, want)
		}
	}
}

func TestParseDurationMinutes(t *testing.T) {
	if got := parseDurationMinutes("PT1H20M"); got == nil || *got != 80 {
		t.Fatalf("PT1H20M=%v", got)
	}
	if got := parseDurationMinutes("PT30M"); got == nil || *got != 30 {
		t.Fatalf("PT30M=%v", got)
	}
	if parseDurationMinutes("") != nil {
		t.Fatal("empty should be nil")
	}
}

func TestCleanTitle(t *testing.T) {
	got := cleanTitle("红烧肉 | 下厨房")
	if got != "红烧肉" {
		t.Fatalf("got %q", got)
	}
	if cleanTitle("番茄炒蛋 - 小红书") != "番茄炒蛋" {
		t.Fatalf("suffix: %q", cleanTitle("番茄炒蛋 - 小红书"))
	}
}

func TestParseHTML_OGAndRecipe(t *testing.T) {
	htmlDoc := `<!doctype html><html><head>
<title>Fallback Title | Site</title>
<meta property="og:title" content="OG 菜名" />
<meta property="og:image" content="https://cdn.example.com/og.jpg" />
<meta property="og:site_name" content="Example" />
<script type="application/ld+json">
{
  "@context": "https://schema.org/",
  "@type": "Recipe",
  "name": "Recipe 菜名",
  "image": ["https://cdn.example.com/recipe.jpg"],
  "recipeIngredient": ["鸡肉 300g", "豆腐 1块", "辣椒"],
  "cookTime": "PT1H20M"
}
</script>
</head><body></body></html>`
	meta := parseHTML(strings.NewReader(htmlDoc))
	if meta.OGTitle != "OG 菜名" {
		t.Fatalf("og title %q", meta.OGTitle)
	}
	if meta.RecipeName != "Recipe 菜名" {
		t.Fatalf("recipe name %q", meta.RecipeName)
	}
	if meta.Ingredients != "鸡肉 300g，豆腐 1块，辣椒" {
		t.Fatalf("ingredients %q", meta.Ingredients)
	}
	if meta.CookMinutes == nil || *meta.CookMinutes != 80 {
		t.Fatalf("cook %v", meta.CookMinutes)
	}
	name := pickName(meta, "")
	if name != "Recipe 菜名" {
		t.Fatalf("pickName %q", name)
	}
	cover := pickCover(meta, "")
	if cover != "https://cdn.example.com/recipe.jpg" {
		t.Fatalf("cover %q", cover)
	}
}

func TestParseHTML_GraphRecipe(t *testing.T) {
	htmlDoc := `<html><head><script type="application/ld+json">
{"@graph":[{"@type":"WebSite","name":"x"},{"@type":["Recipe"],"name":"图菜","recipeIngredient":["盐"]}]}
</script></head></html>`
	meta := parseHTML(strings.NewReader(htmlDoc))
	if meta.RecipeName != "图菜" || meta.Ingredients != "盐" {
		t.Fatalf("%+v", meta)
	}
}

func TestSSRFBlocked(t *testing.T) {
	ctx := context.Background()
	u, err := parseHTTPURL("http://127.0.0.1:8080/")
	if err != nil {
		t.Fatal(err)
	}
	if err := validateFetchURL(ctx, u); err == nil {
		t.Fatal("expected block")
	}
	u, _ = parseHTTPURL("http://192.168.0.1/")
	if err := validateFetchURL(ctx, u); err == nil {
		t.Fatal("expected private block")
	}
	if _, err := parseHTTPURL("file:///etc/passwd"); err != ErrInvalidURL {
		t.Fatalf("file: %v", err)
	}
	if _, err := parseHTTPURL("not a url"); err != ErrInvalidURL {
		t.Fatalf("raw: %v", err)
	}
}

func TestPreview_SSRFPartial(t *testing.T) {
	svc := New()
	got, err := svc.Preview(context.Background(), "http://127.0.0.1:8080/")
	if err != nil {
		t.Fatal(err)
	}
	if !got.Partial {
		t.Fatalf("want partial, got %+v", got)
	}
}

func TestPreview_HTMLViaTestServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><head>
<meta property="og:title" content="宫保鸡丁" />
<meta property="og:image" content="https://cdn.example.com/g.jpg" />
</head></html>`))
	}))
	defer ts.Close()
	svc := NewTest(ts.Client(), "", "")
	got, err := svc.Preview(context.Background(), ts.URL+"/recipe")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "宫保鸡丁" || got.CoverImageURL != "https://cdn.example.com/g.jpg" {
		t.Fatalf("%+v", got)
	}
	if got.Partial {
		t.Fatal("should not be partial")
	}
}

func TestPreview_OEmbedYouTube(t *testing.T) {
	pageURL := "https://www.youtube.com/watch?v=abc"
	mux := http.NewServeMux()
	mux.HandleFunc("/oembed", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("url") != pageURL {
			t.Errorf("url param %q", r.URL.Query().Get("url"))
		}
		_ = json.NewEncoder(w).Encode(map[string]string{
			"title":          "测试做饭视频",
			"thumbnail_url":  "https://i.ytimg.com/vi/abc/hqdefault.jpg",
			"provider_name":  "YouTube",
		})
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	svc := NewTest(ts.Client(), ts.URL+"/oembed", ts.URL+"/oembed")
	got, err := svc.Preview(context.Background(), pageURL)
	if err != nil {
		t.Fatal(err)
	}
	if got.SourcePlatform != "YouTube" {
		t.Fatalf("platform %q", got.SourcePlatform)
	}
	if got.Name != "测试做饭视频" {
		t.Fatalf("name %q", got.Name)
	}
	if got.CoverImageURL != "https://i.ytimg.com/vi/abc/hqdefault.jpg" {
		t.Fatalf("cover %q", got.CoverImageURL)
	}
}
