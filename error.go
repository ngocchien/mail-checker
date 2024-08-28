package mail_checker

import "errors"

var (
	ErrMicrosoftGetAmscCookieError   = errors.New("get amsc cookie fail")
	ErrMicrosoftGetCanaryCookieError = errors.New("get canary cookie fail")
)
