package mail_checker

import (
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type yahooMail struct {
	client *http.Client
}

func (y *yahooMail) Check(email string) (status Status) {
	arrDataEmail := strings.SplitN(email, "@", 2)

	if len(arrDataEmail) != 2 {
		log.Errorf("Invalid email format: %s", email)
		return getStatusById(StatusIdFormatInvalid)
	}

	dataBody, err := y.getBodyData()
	if err != nil {
		log.Errorf("Error fetching body data: %v", err)
		return getStatusById(StatusIdCheckError)
	}

	dataBody.UserId = email
	dataBody.UseridDomain = arrDataEmail[1]

	data, err := query.Values(&dataBody)
	if err != nil {
		log.Errorf("Error encoding query data: %v", err)
		return getStatusById(StatusIdCheckError)
	}

	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, yahooCheckerUrlApi, body)
	if err != nil {
		log.Errorf("Error creating new request to %s: %v", yahooCheckerUrlApi, err)
		return getStatusById(StatusIdCheckError)
	}

	req.Header.Set("Content-Type", `application/x-www-form-urlencoded; charset=UTF-8`)
	req.Header.Set("Cookie", dataBody.Cookie)
	req.Header.Set("X-Requested-With", `XMLHttpRequest`)

	resp, err := y.client.Do(req)
	if err != nil {
		log.Errorf("Error executing request: %v", err)
		return getStatusById(StatusIdCheckError)
	}
	defer resp.Body.Close()
	y.client.CloseIdleConnections()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return getStatusById(StatusIdCheckError)
	}

	var responseData yahooResChecker
	if err = json.Unmarshal(bodyBytes, &responseData); err != nil {
		log.Errorf("Error unmarshaling response JSON: %v", err)
		return getStatusById(StatusIdCheckError)
	}

	if responseData.Errors == nil {
		log.Error("No errors field in response data")
		return getStatusById(StatusIdCheckError)
	}

	for _, er := range responseData.Errors {
		if er.Name == yahooKeyCheckExists {
			switch er.Error {
			case yahooTextDetectUnavailableMail,
				yahooTextDetectNotUnavailableMail,
				yahooTextDetectReservedWordPresentMail:
				return getStatusById(StatusIdLive)
			case yahooTextDetectErrorLengthTooShort,
				yahooTextDetectErrorSomeSpecialCharNotAllow:
				return getStatusById(StatusIdCheckError)
			}
		}
	}
	return getStatusById(StatusIdNotExists)
}

func (y *yahooMail) detectValue(html, name string) (string, error) {
	re := regexp.MustCompile(`(?m)value="(.*?)" name="` + name + `"`)
	matches := re.FindStringSubmatch(html)
	if len(matches) == 0 {
		return "", errors.New("could not detect value for " + name)
	}
	return matches[1], nil
}

func (y *yahooMail) getBodyData() (yahooBodyChecker, error) {
	req, err := http.NewRequest(http.MethodGet, yahooCreateAccountUrl, nil)
	if err != nil {
		log.Errorf("Error creating request to %s: %v", yahooCreateAccountUrl, err)
		return yahooBodyChecker{}, err
	}

	res, err := y.client.Do(req)
	if err != nil {
		log.Errorf("Error executing request to %s: %v", yahooCreateAccountUrl, err)
		return yahooBodyChecker{}, err
	}
	defer res.Body.Close()

	cookies := res.Header.Get("Set-Cookie")
	if cookies == "" {
		err := errors.New("could not detect cookies")
		log.Error(err)
		return yahooBodyChecker{}, err
	}
	arrCookies := strings.Split(cookies, ";")
	dataBody := yahooBodyChecker{Cookie: arrCookies[0]}

	htmlBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return yahooBodyChecker{}, err
	}
	html := string(htmlBytes)

	if dataBody.Acrumb, err = y.detectValue(html, "acrumb"); err != nil {
		return yahooBodyChecker{}, err
	}
	if dataBody.Crumb, err = y.detectValue(html, "crumb"); err != nil {
		return yahooBodyChecker{}, err
	}
	if dataBody.SessionIndex, err = y.detectValue(html, "sessionIndex"); err != nil {
		return yahooBodyChecker{}, err
	}
	if dataBody.Tos0, err = y.detectValue(html, "tos0"); err != nil {
		return yahooBodyChecker{}, err
	}
	if dataBody.SpecId, err = y.detectValue(html, "specId"); err != nil {
		return yahooBodyChecker{}, err
	}

	return dataBody, nil
}
