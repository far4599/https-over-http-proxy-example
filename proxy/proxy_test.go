package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/far4599/https-over-http-proxy-example/proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleTargetPageContent = "ok"

func newTargetServer(t *testing.T, tls bool) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(sampleTargetPageContent))
		assert.NoError(t, err)
	})

	if tls {
		return httptest.NewTLSServer(handler)
	}

	return httptest.NewServer(handler)
}

func addProxyURLToClient(t *testing.T, client *http.Client, proxyServerURL string) {
	proxyUrl, err := url.Parse(proxyServerURL)
	assert.NoError(t, err)
	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	transport.Proxy = http.ProxyURL(proxyUrl)
	client.Transport = transport
}

func TestHTTPProxyHandlerFunc(t *testing.T) {
	testCases := []struct {
		name  string
		isTLS bool
	}{
		{
			name:  "proxy HTTP request",
			isTLS: false,
		},
	}

	t.Parallel()

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ts := newTargetServer(t, tc.isTLS)
			defer ts.Close()

			proxyServer := httptest.NewServer(http.HandlerFunc(proxy.HTTPProxyHandlerFunc))
			defer proxyServer.Close()

			addProxyURLToClient(t, ts.Client(), proxyServer.URL)

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, sampleTargetPageContent, string(body))
		})
	}
}
