// generated-from:77a1fb855ee39fd95d2afc0ffbba7e3b23493186b541278699e654471dc7e95e DO NOT REMOVE, DO UPDATE

package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/http/cookiejar"

	client "github.com/moov-io/ach-web-viewer/pkg/client"
)

func NewTestClient(handler http.Handler) *client.APIClient {
	mockHandler := MockClientHandler{
		handler: handler,
	}

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	mockClient := &http.Client{
		Jar: cookieJar,

		// Mock handler that sends the request to the handler passed in and returns the response without a server
		// middleman.
		Transport: &mockHandler,

		// Disables following redirects for testing.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	config := client.NewConfiguration()
	config.HTTPClient = mockClient
	apiClient := client.NewAPIClient(config)

	return apiClient
}

type MockClientHandler struct {
	handler http.Handler
	ctx     *context.Context
}

func (h *MockClientHandler) RoundTrip(request *http.Request) (*http.Response, error) {
	writer := httptest.NewRecorder()

	h.handler.ServeHTTP(writer, request)
	return writer.Result(), nil
}
