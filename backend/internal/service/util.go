package service

import (
	"math"
	"strings"
	"unicode"
)

func splitNames(s string) []string {
	s = strings.ReplaceAll(s, "，", ",")
	s = strings.ReplaceAll(s, "、", ",")
	s = strings.ReplaceAll(s, ";", ",")
	s = strings.ReplaceAll(s, "；", ",")
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimFunc(p, unicode.IsSpace)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
