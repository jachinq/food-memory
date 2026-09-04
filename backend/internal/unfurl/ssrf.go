package unfurl

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxRedirects = 5
	maxBodyBytes = 512 * 1024
	httpTimeout  = 8 * time.Second
)

var (
	ErrInvalidURL = errors.New("无效的链接")
	errBlocked    = errors.New("不允许访问该地址")
)

func parseHTTPURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, ErrInvalidURL
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return nil, ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, ErrInvalidURL
	}
	return u, nil
}

func blockedHost(host string) bool {
	h := strings.ToLower(strings.TrimSpace(host))
	if h == "" {
		return true
	}
	if strings.Contains(h, ":") {
		if hostOnly, _, err := net.SplitHostPort(h); err == nil {
			h = hostOnly
		}
	}
	h = strings.Trim(h, "[]")
	switch h {
	case "localhost", "localhost.localdomain", "metadata.google.internal":
		return true
	}
	if strings.HasSuffix(h, ".localhost") || strings.HasSuffix(h, ".local") {
		return true
	}
	if ip := net.ParseIP(h); ip != nil {
		return blockedIP(ip)
	}
	return false
}

func blockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}
	return false
}

func resolveBlocked(ctx context.Context, host string) error {
	h := host
	if hostOnly, _, err := net.SplitHostPort(host); err == nil {
		h = hostOnly
	}
	h = strings.Trim(h, "[]")
	if blockedHost(h) {
		return errBlocked
	}
	if ip := net.ParseIP(h); ip != nil {
		if blockedIP(ip) {
			return errBlocked
		}
		return nil
	}
	resolver := net.DefaultResolver
	ips, err := resolver.LookupIPAddr(ctx, h)
	if err != nil {
		return fmt.Errorf("解析主机失败: %w", err)
	}
	if len(ips) == 0 {
		return errBlocked
	}
	for _, addr := range ips {
		if blockedIP(addr.IP) {
			return errBlocked
		}
	}
	return nil
}

func validateFetchURL(ctx context.Context, u *url.URL) error {
	if u == nil {
		return ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	return resolveBlocked(ctx, u.Host)
}

func newSafeClient() *http.Client {
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	transport := &http.Transport{
		Proxy: nil,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			if err := resolveBlocked(ctx, host); err != nil {
				return nil, err
			}
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, err
			}
			var last error
			for _, ipa := range ips {
				if blockedIP(ipa.IP) {
					last = errBlocked
					continue
				}
				conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ipa.IP.String(), port))
				if err == nil {
					return conn, nil
				}
				last = err
			}
			if last == nil {
				last = errBlocked
			}
			return nil, last
		},
		DisableKeepAlives: true,
	}
	return &http.Client{
		Timeout:   httpTimeout,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return errors.New("重定向过多")
			}
			return validateFetchURL(req.Context(), req.URL)
		},
	}
}
