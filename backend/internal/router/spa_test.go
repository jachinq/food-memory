package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func spaEngine(t *testing.T) *gin.Engine {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("INDEX"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "assets", "app.js"), []byte("ASSET"), 0o644); err != nil {
		t.Fatal(err)
	}
	for name, body := range map[string]string{
		"sw.js":                 "SW",
		"registerSW.js":         "REG",
		"workbox-deadbeef.js":   "WB",
		"manifest.webmanifest":  "{}",
		"favicon.svg":           "<svg/>",
		"pwa-192.png":           "PNG192",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	r := gin.New()
	mountSPA(r, dir)
	return r
}

func get(r http.Handler, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestMountSPAServesServiceWorkerWithNoCache(t *testing.T) {
	r := spaEngine(t)
	w := get(r, "/sw.js")
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	if w.Body.String() != "SW" {
		t.Fatalf("body %q", w.Body.String())
	}
	if !strings.Contains(w.Header().Get("Cache-Control"), "no-cache") {
		t.Fatalf("Cache-Control %q", w.Header().Get("Cache-Control"))
	}
}

func TestMountSPAServesWorkboxAndRegisterScriptsWithNoCache(t *testing.T) {
	r := spaEngine(t)
	for _, path := range []string{"/registerSW.js", "/workbox-deadbeef.js"} {
		w := get(r, path)
		if w.Code != http.StatusOK {
			t.Fatalf("%s status %d", path, w.Code)
		}
		if !strings.Contains(w.Header().Get("Cache-Control"), "no-cache") {
			t.Fatalf("%s Cache-Control %q", path, w.Header().Get("Cache-Control"))
		}
	}
}

func TestMountSPAServesManifestAndIcons(t *testing.T) {
	r := spaEngine(t)
	man := get(r, "/manifest.webmanifest")
	if man.Code != http.StatusOK || man.Body.String() != "{}" {
		t.Fatalf("manifest %d %q", man.Code, man.Body.String())
	}
	icon := get(r, "/pwa-192.png")
	if icon.Code != http.StatusOK || icon.Body.String() != "PNG192" {
		t.Fatalf("icon %d %q", icon.Code, icon.Body.String())
	}
}

func TestMountSPAUnknownPathReturnsIndex(t *testing.T) {
	r := spaEngine(t)
	w := get(r, "/me")
	if w.Code != http.StatusOK || w.Body.String() != "INDEX" {
		t.Fatalf("got %d %q", w.Code, w.Body.String())
	}
}

func TestMountSPAKeepsHashedAssets(t *testing.T) {
	r := spaEngine(t)
	w := get(r, "/assets/app.js")
	if w.Code != http.StatusOK || w.Body.String() != "ASSET" {
		t.Fatalf("got %d %q", w.Code, w.Body.String())
	}
}

func TestMountSPAAPIPrefixIsNotIndex(t *testing.T) {
	r := spaEngine(t)
	w := get(r, "/api/missing")
	if w.Code != http.StatusNotFound {
		t.Fatalf("status %d", w.Code)
	}
	if strings.Contains(w.Body.String(), "INDEX") {
		t.Fatalf("api 404 served index: %q", w.Body.String())
	}
}
