# Mail Availability Checker

This Go package allows you to check the availability of Microsoft email addresses (e.g., Outlook, Hotmail). It performs an HTTP request to Microsoft's API to determine if an email address is already in use.

## Features

- Check if a Microsoft email address is available.
- Supports proxy configuration.
- Comprehensive error handling and logging.
- Unit tested with >90% coverage.

## Installation

To use this package, you need to have Go installed. Then, you can install the package by running:

```bash
go get github.com/ngocchien/mail_checker
```

## Usage

### Import the Package

```go
import "github.com/ngocchien/mail_checker"
```

### Create a New Checker

To create a new instance of the Microsoft Mail Checker:

```go
checker := mail_checker.New(mail_checker.MailKindMicrosoft, mail_checker.Proxy{})
```

- **MailKindMicrosoft**: This constant represents the Microsoft mail kind.
- **Proxy**: (Optional) If you need to use a proxy, pass a `Proxy` struct with the necessary fields (Host, Schema, User, Password). Otherwise, pass an empty `Proxy{}`.

### Check Email Availability

To check the availability of an email address:

```go
status := checker.Check("email@example.com")
fmt.Printf("Status: %d, Name: %s, Message: %s\n", status.Id, status.Name, status.Message)
```

- **email@example.com**: Replace this with the email address you want to check.
- The `Check` method returns a `Status` struct that contains:
    - `Id`: Status ID (e.g., `StatusIdLive`, `StatusIdNotExists`).
    - `Name`: Status name (e.g., "Live", "Not exists").
    - `Message`: A custom message if any.
    - `Data`: Any additional data.

### Example

```go
package main

import (
	"fmt"
	"github.com/ngocchien/mail_checker"
)

func main() {
	proxy := mail_checker.Proxy{
		Host: "127.0.0.1:8080",
	}

	checker := mail_checker.New(mail_checker.MailKindMicrosoft, proxy)
	status := checker.Check("test@example.com")

	fmt.Printf("Status: %d, Name: %s, Message: %s\n", status.Id, status.Name, status.Message)
}
```

### Proxy Configuration

If you need to use a proxy, create a `Proxy` struct and pass it to the `New` function:

```go
proxy := mail_checker.Proxy{
	Host:     "127.0.0.1:8080",
	Schema:   "http",
	User:     "proxyUser",
	Password: "proxyPassword",
}

checker := mail_checker.New(mail_checker.MailKindMicrosoft, proxy)
```

### Error Handling

The package uses `logrus` for logging errors. Make sure to configure `logrus` according to your application's needs.

### Unit Tests

This package is thoroughly tested with more than 90% code coverage. To run the tests, use:

```bash
go test ./... -cover
```

### Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---