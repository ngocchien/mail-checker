package mail_checker

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Checker interface {
	Check(email string) (status Status)
}

func makeHttpClient(proxy Proxy) *http.Client {
	c := http.Client{
		Timeout: httpClientTimeoutDefault,
	}
	if proxy.Host != "" {
		p := &url.URL{
			Host: proxy.Host,
		}
		if proxy.User != "" && proxy.Password != "" {
			p.User = url.UserPassword(proxy.User, proxy.Password)
		}
		c.Transport = &http.Transport{
			Proxy: http.ProxyURL(p),
		}
	}
	return &c
}

func getStatusById(id StatusId) (status Status) {
	switch id {
	case StatusIdLive:
		status = Status{
			Id:   id,
			Name: StatusNameLive,
		}
	case StatusIdNotExists:
		status = Status{
			Id:   id,
			Name: StatusNameNotExists,
		}
	case StatusIdDisable:
		status = Status{
			Id:   id,
			Name: StatusNameDisable,
		}
	case StatusIdVerPhone:
		status = Status{
			Id:   id,
			Name: StatusNameVerPhone,
		}
	case StatusIdCheckError:
		status = Status{
			Id:   id,
			Name: StatusNameCheckError,
		}
	case StatusIdFormatInvalid:
		status = Status{
			Id:   id,
			Name: StatusNameFormatInvalid,
		}
	}
	return status
}

func New(mailKind MailKind, proxy Proxy) Checker {
	client := makeHttpClient(proxy)
	switch mailKind {
	case MailKindMicrosoft:
		return &microsoftMail{
			client: client,
		}
	case MailKindYahoo:
		return &yahooMail{
			client: client,
		}
	default:
		log.Errorf("The mail kind input invalid")
	}
	return nil
}
