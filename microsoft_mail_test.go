package mail_checker

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Test the getAmscCookie function for success case
func TestGetAmscCookie_Success(t *testing.T) {
	res := &http.Response{
		Header: http.Header{
			"Set-Cookie": {"amsc=testCookie; path=/;"},
		},
	}
	err, amscCookie := getAmscCookie(res)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amscCookie != "testCookie" {
		t.Fatalf("expected 'testCookie', got %v", amscCookie)
	}
}

// Test the getAmscCookie function for error case
func TestGetAmscCookie_Error(t *testing.T) {
	res := &http.Response{
		Header: http.Header{
			"Set-Cookie": {""},
		},
	}
	err, _ := getAmscCookie(res)
	if err != ErrMicrosoftGetAmscCookieError {
		t.Fatalf("expected ErrMicrosoftGetAmscCookieError, got %v", err)
	}
}

// Test the getCanaryCookie function for success case
func TestGetCanaryCookie_Success(t *testing.T) {
	html := `var ServerData={"apiCanary":"testCanary"};`
	mailChecker := &microsoftMail{}
	err1, canary := mailChecker.getCanaryCookie(html)
	if err1 != nil {
		t.Fatalf("expected no error, got %v", err1)
	}
	//log.Infof("TestGetCanaryCookie_Success canary: %+v", canary)
	if canary != "testCanary" {
		t.Fatalf("expected 'testCanary', got %v", canary)
	}
}

// Test the getCanaryCookie function for JSON parsing error
func TestGetCanaryCookie_JsonParseError(t *testing.T) {
	html := `var ServerData={"apiCanary":invalidJson};`
	mailChecker := &microsoftMail{}
	err, _ := mailChecker.getCanaryCookie(html)
	if err == nil {
		t.Fatalf("expected JSON parsing error, got nil")
	}
	if err != ErrMicrosoftGetCanaryCookieError {
		t.Fatalf("expected ErrMicrosoftGetCanaryCookieError, got %v", err)
	}
}

// Test the getCanaryCookie function for missing canary
func TestGetCanaryCookie_MissingCanary(t *testing.T) {
	html := `var ServerData={};`
	mailChecker := &microsoftMail{}
	err, _ := mailChecker.getCanaryCookie(html)
	if err != ErrMicrosoftGetCanaryCookieError {
		t.Fatalf("expected ErrMicrosoftGetCanaryCookieError, got %v", err)
	}
}

// Test the getAmscAndCanaryCookie function for success case
func TestGetAmscAndCanaryCookie_Success(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == hotmailUrlSignup {
					return &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Set-Cookie": {"amsc=testCookie; path=/;"},
						},
						Body: io.NopCloser(strings.NewReader(`var ServerData={"apiCanary":"testCanary"};`)),
					}, nil
				}
				return nil, nil
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	err, amscCookie, canary := mailChecker.getAmscAndCanaryCookie()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if amscCookie != "testCookie" {
		t.Fatalf("expected 'testCookie', got %v", amscCookie)
	}
	if canary != "testCanary" {
		t.Fatalf("expected 'testCanary', got %v", canary)
	}
}

// Test the getAmscAndCanaryCookie function for HTTP request error
func TestGetAmscAndCanaryCookie_HttpRequestError(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("mock error")
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	err, _, _ := mailChecker.getAmscAndCanaryCookie()
	if err == nil {
		t.Fatalf("expected HTTP request error, got nil")
	}
}

// Test the Check function for a non-available email
func TestCheck_EmailNotAvailable(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == hotmailUrlSignup {
					return &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Set-Cookie": {"amsc=testCookie; path=/;"},
						},
						Body: io.NopCloser(strings.NewReader(`var ServerData={"apiCanary":"testCanary"};`)),
					}, nil
				} else if req.URL.String() == hotmailUrlCheckAvailable {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"isAvailable":true}`)),
					}, nil
				}
				return nil, nil
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	status := mailChecker.Check("test@example.com")
	if status.Id != StatusIdNotExists {
		t.Fatalf("expected StatusIdNotExists, got %v", status.Id)
	}
}

// Test the Check function for an available email
func TestCheck_EmailAvailable(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == hotmailUrlSignup {
					return &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Set-Cookie": {"amsc=testCookie; path=/;"},
						},
						Body: io.NopCloser(strings.NewReader(`var ServerData={"apiCanary":"testCanary"};`)),
					}, nil
				} else if req.URL.String() == hotmailUrlCheckAvailable {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"isAvailable":false}`)),
					}, nil
				}
				return nil, nil
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	status := mailChecker.Check("test@example.com")
	if status.Id != StatusIdLive {
		t.Fatalf("expected StatusIdLive, got %v", status.Id)
	}
}

// Test the Check function for request error
func TestCheck_RequestError(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("mock error")
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	status := mailChecker.Check("test@example.com")
	if status.Id != StatusIdCheckError {
		t.Fatalf("expected StatusIdCheckError, got %v", status.Id)
	}
}

// Test the Check function for malformed JSON in response
func TestCheck_MalformedJsonResponse(t *testing.T) {
	client := &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				if req.URL.String() == hotmailUrlSignup {
					return &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Set-Cookie": {"amsc=testCookie; path=/;"},
						},
						Body: io.NopCloser(strings.NewReader(`var ServerData={"apiCanary":"testCanary"};`)),
					}, nil
				} else if req.URL.String() == hotmailUrlCheckAvailable {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"isAvailable":invalidJson}`)),
					}, nil
				}
				return nil, nil
			},
		},
	}

	mailChecker := &microsoftMail{client: client}
	status := mailChecker.Check("test@example.com")
	if status.Id != StatusIdCheckError {
		t.Fatalf("expected StatusIdCheckError, got %v", status.Id)
	}
}
