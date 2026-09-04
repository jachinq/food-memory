package unfurl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type oembedResp struct {
	Title         string `json:"title"`
	ThumbnailURL  string `json:"thumbnail_url"`
	ProviderName  string `json:"provider_name"`
}

const (
	defaultYouTubeOEmbed = "https://www.youtube.com/oembed"
	defaultTikTokOEmbed  = "https://www.tiktok.com/oembed"
)

func (s *Service) oembedEndpoint(u *url.URL) string {
	if isYouTubeHost(u.Host) {
		return s.youtubeOEmbed
	}
	if isTikTokHost(u.Host) {
		return s.tiktokOEmbed
	}
	return ""
}

func (s *Service) fetchOEmbed(ctx context.Context, pageURL *url.URL) (oembedResp, error) {
	endpoint := s.oembedEndpoint(pageURL)
	if endpoint == "" {
		return oembedResp{}, fmt.Errorf("no oembed")
	}
	q := url.Values{}
	q.Set("format", "json")
	q.Set("url", pageURL.String())
	reqURL := endpoint
	if strings.Contains(endpoint, "?") {
		reqURL += "&" + q.Encode()
	} else {
		reqURL += "?" + q.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return oembedResp{}, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		return oembedResp{}, err
	}
	defer res.Body.Close()
	body, err := readLimited(res.Body)
	if err != nil {
		return oembedResp{}, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return oembedResp{}, fmt.Errorf("oembed status %d", res.StatusCode)
	}
	var out oembedResp
	if err := json.Unmarshal(body, &out); err != nil {
		return oembedResp{}, err
	}
	return out, nil
}
