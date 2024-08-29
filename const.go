package mail_checker

import "time"

const (
	StatusIdLive          StatusId = 1
	StatusIdNotExists     StatusId = 2
	StatusIdDisable       StatusId = 3
	StatusIdVerPhone      StatusId = 4
	StatusIdCheckError    StatusId = 5
	StatusIdFormatInvalid StatusId = 6

	StatusNameLive          StatusName = "Live"
	StatusNameNotExists     StatusName = "Not exists"
	StatusNameDisable       StatusName = "Disable"
	StatusNameVerPhone      StatusName = "Ver phone"
	StatusNameCheckError    StatusName = "Check error"
	StatusNameFormatInvalid StatusName = "Format Invalid"
)

const (
	MailKindMicrosoft        MailKind = "microsoft"
	MailKindGoogle           MailKind = "google"
	MailKindYahoo            MailKind = "yahoo"
	dialProtocol                      = "tcp"
	hotmailUrlSignup                  = "https://signup.live.com/signup"
	hotmailUrlCheckAvailable          = "https://signup.live.com/API/CheckAvailableSigninNames"

	yahooCreateAccountUrl                       = "https://login.yahoo.com/account/create"
	yahooCheckerUrlApi                          = "https://login.yahoo.com/account/module/create?validateField=userId"
	yahooKeyCheckExists                         = "userId"
	yahooTextDetectUnavailableMail              = "IDENTIFIER_EXISTS"
	yahooTextDetectNotUnavailableMail           = "IDENTIFIER_NOT_AVAILABLE"
	yahooTextDetectReservedWordPresentMail      = "RESERVED_WORD_PRESENT"
	yahooTextDetectErrorLengthTooShort          = "LENGTH_TOO_SHORT"
	yahooTextDetectErrorSomeSpecialCharNotAllow = "SOME_SPECIAL_CHARACTERS_NOT_ALLOWED"

	httpClientTimeoutDefault = 5 * time.Second
)
