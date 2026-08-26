package regression

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/karst-pulse/karst-pulse/internal/transport"
)

func TestBug09_HTTPErrorWrappingKeepsStatus(t *testing.T) {
    wrapped := fmt.Errorf("route lookup: %w", transport.NotFound("spring missing"))
    if got := transport.StatusOf(wrapped); got != http.StatusNotFound {
        t.Fatalf("status=%d, want %d", got, http.StatusNotFound)
    }
    response := httptest.NewRecorder()
    transport.WriteError(response, wrapped)
    if response.Code != http.StatusNotFound {
        t.Fatalf("response status=%d, want %d", response.Code, http.StatusNotFound)
    }
}
