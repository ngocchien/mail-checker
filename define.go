package mail_checker

type (
	StatusId   int
	StatusName string
	MailKind   string

	Proxy struct {
		Host     string
		Schema   string
		User     string
		Password string
	}

	Status struct {
		Id   StatusId   `json:"id"`
		Name StatusName `json:"name"`
	}

	microsoftMailResCanary struct {
		ApiCanary string `json:"apiCanary"`
	}
	microsoftMailResResGetEmailAvailable struct {
		IsAvailable bool `json:"isAvailable"`
	}

	yahooBodyChecker struct {
		SpecId        string `url:"specId"`
		CacheStored   string `url:"cacheStored"`
		Crumb         string `url:"crumb"`
		Acrumb        string `url:"acrumb"`
		GoogleIdToken string `url:"googleIdToken"`
		AuthCode      string `url:"authCode"`
		AttrSetIndex  string `url:"attrSetIndex"`
		MultiDomain   string `url:"multiDomain"`
		FirstName     string `url:"firstName"`
		LastName      string `url:"lastName"`
		UseridDomain  string `url:"userid-domain"`
		UserId        string `url:"userId"`
		Password      string `url:"password"`
		Signup        string `url:"signup"`
		SessionIndex  string `url:"sessionIndex"`
		Tos0          string `url:"tos0"`
		Cookie        string `url:"-"`
	}

	yahooResChecker struct {
		Errors []struct {
			Name  string `json:"name"`
			Error string `json:"error"`
		} `json:"errors"`
	}
)
