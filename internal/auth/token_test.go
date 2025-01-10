package auth

import(
	"testing"
	"net/http"
)

func TestGetBearerToken(t *testing.T) {
	t.Run("valid header", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer ThisIsAToken")
		want := "ThisIsAToken"
			
		bearerToken, err := GetBearerToken(headers)
		if err != nil {
			t.Errorf("Actual: bearerToken = %s, Wanted: bearerToken = %s and error: %v", bearerToken, want, err)
		}
	})
	t.Run("empty header", func(t *testing.T) {
		headers := http.Header{}
			
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("Expected error due to empty header, got nil")
		}
	})
	t.Run("invalid authorization header", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "BAERER ThisIsAToken")
			
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("Expected error due to wrong authorization header content, got nil")
		}
	})
	t.Run("empty token in header", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Authorization", "Bearer ")
			
		_, err := GetBearerToken(headers)
		if err == nil {
			t.Error("Expected error due to empty token, got nil")
		}
	})
}