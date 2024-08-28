package mail_checker

import (
	"net/http"
	"testing"
)

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// Test the makeHttpClient function
func TestMakeHttpClient(t *testing.T) {
	proxy := Proxy{
		Host: "127.0.0.1:8080",
	}
	client := makeHttpClient(proxy)
	if client.Timeout != httpClientTimeoutDefault {
		t.Errorf("expected timeout %v, got %v", httpClientTimeoutDefault, client.Timeout)
	}
}

// Test getStatusById function
func TestGetStatusById(t *testing.T) {
	status := getStatusById(StatusIdLive)
	if status.Id != StatusIdLive {
		t.Errorf("expected %v, got %v", StatusIdLive, status.Id)
	}
	if status.Name != StatusNameLive {
		t.Errorf("expected %v, got %v", StatusNameLive, status.Name)
	}
}

// Test New function
func TestNew(t *testing.T) {
	checker := New(MailKindMicrosoft, Proxy{})
	if checker == nil {
		t.Errorf("expected a non-nil checker")
	}
}
