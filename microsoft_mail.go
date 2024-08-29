package mail_checker

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type microsoftMail struct {
	client *http.Client
}

func (h *microsoftMail) Check(email string) (status Status) {
	err, amscCookie, canary := h.getAmscAndCanaryCookie()
	if err != nil {
		log.Errorf("[MicrosoftMail] - [Check] - %s", err.Error())
		return getStatusById(StatusIdCheckError)
	}

	var bodyReq = map[string]interface{}{
		"signInName":         email,
		"includeSuggestions": true,
	}
	var body, _ = json.Marshal(bodyReq)

	r, err := http.NewRequest(http.MethodPost, hotmailUrlCheckAvailable, bytes.NewBuffer(body))
	if err != nil {
		return getStatusById(StatusIdCheckError)
	}
	r.Header.Set("canary", canary)
	r.Header.Set("content-type", "application/json")
	r.Header.Set("cookie", `amsc=`+amscCookie+`;`)
	res, err := h.client.Do(r)
	h.client.CloseIdleConnections()

	if err != nil {
		log.Errorf("Exec request: %+v", err)
		return getStatusById(StatusIdCheckError)
	}

	defer res.Body.Close()
	bodyText, _ := io.ReadAll(res.Body)
	jsonString := string(bodyText)
	if !strings.Contains(jsonString, `isAvailable`) {
		log.Errorf("[MicrosoftMail] - [Check] - The isAvailable field does not exsist in the response")
		return getStatusById(StatusIdCheckError)
	}

	var checkerResponse microsoftMailResResGetEmailAvailable
	err = json.Unmarshal(bodyText, &checkerResponse)
	if err != nil {
		log.Errorf("[MicrosoftMail] - [Check] - Parser JsonBody error: %+v", err)
		return getStatusById(StatusIdCheckError)
	}

	if checkerResponse.IsAvailable {
		return getStatusById(StatusIdNotExists)
	}
	return getStatusById(StatusIdLive)
}

func (h *microsoftMail) getAmscCookie(res *http.Response) (err error, amscCookie string) {
	cookies := res.Header.Get("Set-Cookie")
	re := regexp.MustCompile(`(?s)amsc=(.*?);`)
	cookiesMatches := re.FindStringSubmatch(cookies)
	if len(cookiesMatches) == 0 {
		return ErrMicrosoftGetAmscCookieError, amscCookie
	}
	return nil, cookiesMatches[1]
}

func (h *microsoftMail) getCanaryCookie(html string) (err error, canary string) {
	re := regexp.MustCompile(`(?s)var ServerData=(.*?);`)
	canaryCookiesMatches := re.FindStringSubmatch(html)

	if len(canaryCookiesMatches) == 0 {
		return ErrMicrosoftGetCanaryCookieError, canary
	}

	var dataBody microsoftMailResCanary
	err = json.Unmarshal([]byte(canaryCookiesMatches[1]), &dataBody)

	if err != nil {
		return ErrMicrosoftGetCanaryCookieError, canary
	}

	if dataBody.ApiCanary == "" {
		return ErrMicrosoftGetCanaryCookieError, canary
	}

	return nil, dataBody.ApiCanary
}

func (h *microsoftMail) getAmscAndCanaryCookie() (err error, amscCookie string, canary string) {
	r, err := http.NewRequest(http.MethodGet, hotmailUrlSignup, nil)
	if err != nil {
		return err, amscCookie, amscCookie
	}
	res, err := h.client.Do(r)
	h.client.CloseIdleConnections()
	if err != nil {
		return err, amscCookie, canary
	}
	defer res.Body.Close()

	err, amscCookie = h.getAmscCookie(res)
	if err != nil {
		return err, amscCookie, canary
	}

	htmlByte, _ := io.ReadAll(res.Body)
	err, canary = h.getCanaryCookie(string(htmlByte))
	if err != nil {
		return fmt.Errorf("failed to get canary cookie: %w", err), amscCookie, canary
	}
	return err, amscCookie, canary
}
