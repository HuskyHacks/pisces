package piscestest

import (
	"context"
	"embed"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/mjc-gh/pisces/internal/browser"
)

//go:embed testdata/*
var testFS embed.FS

type handler struct {
	dir string
}

func NewTestWebServer(dir string) *httptest.Server {
	server := httptest.NewServer(handler{dir})

	return server
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the file path from the request
	path := r.URL.Path
	fullPath := filepath.Join("testdata", h.dir, path)

	file, err := testFS.Open(fullPath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	} else if _, err := file.Stat(); errors.Is(err, os.ErrNotExist) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Serve the file
	http.ServeFileFS(w, r, testFS, fullPath)
}

func NewTestContext() context.Context {
	bCtx, _ := browserTestContext()

	return bCtx
}

func browserTestContext() (context.Context, context.CancelFunc) {
	ctx := context.TODO()
	remoteUrl, useRemote := os.LookupEnv("PISCES_CHROMEDP_REMOTE_URL")
	if useRemote {
		return browser.StartRemote(ctx, remoteUrl)
	}

	return browser.StartLocal(ctx)
}
