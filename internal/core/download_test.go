package core

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloadBookmark(t *testing.T) {
	t.Run("successful download", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body>hello</body></html>"))
		}))
		defer server.Close()

		body, contentType, err := DownloadBookmark(server.URL)
		require.NoError(t, err)
		require.NotNil(t, body)
		defer body.Close()
		require.Equal(t, "text/html", contentType)
	})

	t.Run("returns error on HTTP 403", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Access Denied"))
		}))
		defer server.Close()

		body, _, err := DownloadBookmark(server.URL)
		require.Error(t, err)
		require.Nil(t, body)
		require.Contains(t, err.Error(), "HTTP 403")
	})

	t.Run("returns error on HTTP 404", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		body, _, err := DownloadBookmark(server.URL)
		require.Error(t, err)
		require.Nil(t, body)
		require.Contains(t, err.Error(), "HTTP 404")
	})

	t.Run("returns error on Cloudflare challenge", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("cf-mitigated", "challenge")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html>Cloudflare challenge page</html>"))
		}))
		defer server.Close()

		body, _, err := DownloadBookmark(server.URL)
		require.Error(t, err)
		require.Nil(t, body)
		require.Contains(t, err.Error(), "HTTP 200")
	})

	t.Run("returns error on HTTP 500", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		body, _, err := DownloadBookmark(server.URL)
		require.Error(t, err)
		require.Nil(t, body)
		require.Contains(t, err.Error(), "HTTP 500")
	})
}
