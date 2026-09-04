package unfurl

import (
	"context"
	"io"
	"net/http"
	"strings"
)

const userAgent = "ShiyiLinkPreview/1.0"

type Result struct {
	URL             string
	SourcePlatform  string
	Name            string
	CoverImageURL   string
	MainIngredients string
	CookTimeMinutes *int
	Partial         bool
}

type Service struct {
	client        *http.Client
	youtubeOEmbed string
	tiktokOEmbed  string
	skipSSRF      bool
}

func New() *Service {
	return &Service{
		client:        newSafeClient(),
		youtubeOEmbed: defaultYouTubeOEmbed,
		tiktokOEmbed:  defaultTikTokOEmbed,
	}
}

func NewTest(client *http.Client, youtubeOEmbed, tiktokOEmbed string) *Service {
	if client == nil {
		client = http.DefaultClient
	}
	return &Service{
		client:        client,
		youtubeOEmbed: youtubeOEmbed,
		tiktokOEmbed:  tiktokOEmbed,
		skipSSRF:      true,
	}
}

func (s *Service) Preview(ctx context.Context, rawURL string) (Result, error) {
	u, err := parseHTTPURL(rawURL)
	if err != nil {
		return Result{}, ErrInvalidURL
	}
	out := Result{URL: u.String(), SourcePlatform: PlatformFromHost(u.Host)}
	if !s.skipSSRF {
		if err := validateFetchURL(ctx, u); err != nil {
			if err == ErrInvalidURL {
				return Result{}, err
			}
			out.Partial = true
			return out, nil
		}
	}

	fetched := false
	var meta pageMeta
	var oe oembedResp

	if s.oembedEndpoint(u) != "" {
		if got, err := s.fetchOEmbed(ctx, u); err == nil {
			oe = got
			fetched = true
		}
	} else if body, err := s.fetchPage(ctx, u.String()); err == nil {
		meta = parseHTML(strings.NewReader(string(body)))
		fetched = true
	}

	out.Name = pickName(meta, oe.Title)
	out.CoverImageURL = pickCover(meta, oe.ThumbnailURL)
	out.MainIngredients = meta.Ingredients
	out.CookTimeMinutes = meta.CookMinutes
	out.SourcePlatform = pickPlatform(out.SourcePlatform, meta.OGSiteName, oe.ProviderName)
	out.Partial = !fetched || (out.Name == "" && out.CoverImageURL == "" && out.MainIngredients == "" && out.CookTimeMinutes == nil)
	return out, nil
}

func (s *Service) fetchPage(ctx context.Context, raw string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.8")
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(io.LimitReader(res.Body, maxBodyBytes))
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, io.EOF
	}
	return body, nil
}
