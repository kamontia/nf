package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppNotifier_Notify(t *testing.T) {
	testCases := []struct {
		name          string
		apiToken      string
		expectAuthHdr bool
	}{
		{
			name:          "without auth token",
			apiToken:      "",
			expectAuthHdr: false,
		},
		{
			name:          "with auth token",
			apiToken:      "test-token-123",
			expectAuthHdr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 1. Check method
				assert.Equal(t, "POST", r.Method, "Expected POST request")

				// 2. Check auth header
				authHeader := r.Header.Get("Authorization")
				if tc.expectAuthHdr {
					expectedHdr := "Bearer " + tc.apiToken
					assert.Equal(t, expectedHdr, authHeader, "Authorization header is incorrect")
				} else {
					assert.Empty(t, authHeader, "Expected no Authorization header")
				}

				// 3. Check body
				bodyBytes, err := io.ReadAll(r.Body)
				require.NoError(t, err, "Failed to read request body")
				defer r.Body.Close()

				var payload appPayload
				err = json.Unmarshal(bodyBytes, &payload)
				require.NoError(t, err, "Failed to unmarshal request body")

				assert.Equal(t, "Test Title", payload.Title)
				assert.Equal(t, "Test Message", payload.Message)

				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			// Create the notifier pointing to the mock server
			notifier := NewAppNotifier(server.URL, tc.apiToken)

			// Call Notify
			err := notifier.Notify("Test Title", "Test Message")

			// Assert no error was returned
			assert.NoError(t, err, "Notify returned an unexpected error")
		})
	}
}
