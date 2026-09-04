package unfurl

import (
	"net"
	"strings"
)

var hostPlatforms = []struct {
	suffix   string
	platform string
}{
	{suffix: "xiaohongshu.com", platform: "小红书"},
	{suffix: "xhslink.com", platform: "小红书"},
	{suffix: "douyin.com", platform: "抖音"},
	{suffix: "iesdouyin.com", platform: "抖音"},
	{suffix: "bilibili.com", platform: "B站"},
	{suffix: "b23.tv", platform: "B站"},
	{suffix: "weibo.com", platform: "微博"},
	{suffix: "weibo.cn", platform: "微博"},
	{suffix: "mp.weixin.qq.com", platform: "公众号"},
	{suffix: "youtube.com", platform: "YouTube"},
	{suffix: "youtu.be", platform: "YouTube"},
	{suffix: "tiktok.com", platform: "TikTok"},
	{suffix: "xiachufang.com", platform: "下厨房"},
}

func hostnameOf(host string) string {
	h := strings.ToLower(strings.TrimSpace(host))
	if h == "" {
		return ""
	}
	if hostOnly, _, err := net.SplitHostPort(h); err == nil {
		h = hostOnly
	}
	return strings.Trim(h, "[]")
}

func PlatformFromHost(host string) string {
	h := hostnameOf(host)
	if h == "" {
		return ""
	}
	for _, item := range hostPlatforms {
		if h == item.suffix || strings.HasSuffix(h, "."+item.suffix) {
			return item.platform
		}
	}
	return ""
}

func isYouTubeHost(host string) bool {
	h := hostnameOf(host)
	return h == "youtube.com" || strings.HasSuffix(h, ".youtube.com") || h == "youtu.be"
}

func isTikTokHost(host string) bool {
	h := hostnameOf(host)
	return h == "tiktok.com" || strings.HasSuffix(h, ".tiktok.com")
}
