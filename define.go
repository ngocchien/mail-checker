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
)
