package mail_checker

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func newMockClient(roundTripFunc func(req *http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: &mockTransport{RoundTripFunc: roundTripFunc},
	}
}

// Test detectValue
func TestDetectValue(t *testing.T) {
	y := yahooMail{}
	html := `<input type="hidden" value="testValue" name="acrumb">`
	value, err := y.detectValue(html, "acrumb")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "testValue" {
		t.Fatalf("expected 'testValue', got %v", value)
	}

	// Test case where value is not found
	_, err = y.detectValue(html, "nonexistent")
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestGetBodyDataSuccess(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		cookie := `testCookie`
		html := `<input type="hidden" value="acrumb" name="acrumb">
                 <input type="hidden" value="crumb" name="crumb">
                 <input type="hidden" value="sessionIndex" name="sessionIndex">
                 <input type="hidden" value="tos0" name="tos0">
                 <input type="hidden" value="specId" name="specId">`
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Set-Cookie": {cookie}},
			Body:       io.NopCloser(strings.NewReader(html)),
		}, nil
	})
	y := yahooMail{client: client}

	bodyData, err := y.getBodyData()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if bodyData.Cookie != "testCookie" {
		t.Fatalf("expected 'testCookie', got %v", bodyData.Cookie)
	}
	if bodyData.Acrumb != "acrumb" {
		t.Fatalf("expected 'acrumb', got %v", bodyData.Acrumb)
	}
	if bodyData.Crumb != "crumb" {
		t.Fatalf("expected 'crumb', got %v", bodyData.Crumb)
	}
	if bodyData.SessionIndex != "sessionIndex" {
		t.Fatalf("expected 'sessionIndex', got %v", bodyData.SessionIndex)
	}
	if bodyData.Tos0 != "tos0" {
		t.Fatalf("expected 'tos0', got %v", bodyData.Tos0)
	}
	if bodyData.SpecId != "specId" {
		t.Fatalf("expected 'specId', got %v", bodyData.SpecId)
	}
}

// Test getBodyData error cases
func TestGetBodyDataErrors(t *testing.T) {
	// Error in HTTP request
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("mock error")
	})
	y := yahooMail{client: client}

	_, err := y.getBodyData()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Missing cookie
	client = newMockClient(func(req *http.Request) (*http.Response, error) {
		html := `<input type="hidden" value="acrumb" name="acrumb">`
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{}, // No Set-Cookie header
			Body:       io.NopCloser(strings.NewReader(html)),
		}, nil
	})
	y.client = client

	_, err = y.getBodyData()
	if err == nil || err.Error() != "could not detect cookies" {
		t.Fatalf("expected cookie detection error, got %v", err)
	}

	// Error in detecting Acrumb
	client = newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Set-Cookie": {"Set-Cookie: testCookie;"}},
			Body:       io.NopCloser(strings.NewReader("")),
		}, nil
	})
	y.client = client

	_, err = y.getBodyData()
	if err == nil || !strings.Contains(err.Error(), "could not detect value for acrumb") {
		t.Fatalf("expected Acrumb detection error, got %v", err)
	}
}

// Test Check method
func TestCheck(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() == yahooCreateAccountUrl {
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Set-Cookie": {"testCookie"}},
				Body: io.NopCloser(strings.NewReader(`<input type="hidden" value="acrumb" name="acrumb">
                                                       <input type="hidden" value="crumb" name="crumb">
                                                       <input type="hidden" value="sessionIndex" name="sessionIndex">
                                                       <input type="hidden" value="tos0" name="tos0">
                                                       <input type="hidden" value="specId" name="specId">`)),
			}, nil
		} else if req.URL.String() == yahooCheckerUrlApi {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(strings.NewReader(`{
					"errors": [{"name": "userId", "error": "IDENTIFIER_EXISTS"}]
				}`)),
			}, nil
		}
		return nil, errors.New("unexpected URL")
	})
	y := yahooMail{client: client}

	status := y.Check("test@yahoo.com")
	if status.Id != StatusIdLive {
		t.Fatalf("expected StatusIdLive, got %v", status.Id)
	}

	// Test error in request
	client = newMockClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("mock error")
	})
	y.client = client

	status = y.Check("test@yahoo.com")
	if status.Id != StatusIdCheckError {
		t.Fatalf("expected StatusIdCheckError, got %v", status.Id)
	}

	// Test invalid email format
	status = y.Check("invalid-email-format")
	if status.Id != StatusIdFormatInvalid {
		t.Fatalf("expected StatusIdFormatInvalid, got %v", status.Id)
	}

	// Test case with no errors in response
	client = newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(`{"errors": null}`)),
		}, nil
	})
	y.client = client

	status = y.Check("test@yahoo.com")
	if status.Id != StatusIdCheckError {
		t.Fatalf("expected StatusIdCheckError, got %v", status.Id)
	}

	// Test case with specific error types leading to CheckError status
	client = newMockClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(strings.NewReader(`{
				"errors": [{"name": "check_exists", "error": "length_too_short"}]
			}`)),
		}, nil
	})
	y.client = client

	status = y.Check("test@yahoo.com")
	if status.Id != StatusIdCheckError {
		t.Fatalf("expected StatusIdCheckError, got %v", status.Id)
	}
}
