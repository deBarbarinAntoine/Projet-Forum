package utils

import (
	"Projet-Forum/internal/models"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net"
	"net/smtp"
	"os"
)

// configFile is the config's file absolute path.
var configFile = Path + "config/config.json"

// generateConfirmationID
//
//	@Description: generates a random confirmation ID.
//	@return string
func generateConfirmationID() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// generateMessageID
//
//	@Description: generates a message ID to send a mail.
//	@param domain
//	@return string
func generateMessageID(domain string) string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("<%s@%s>", base64.StdEncoding.EncodeToString(b), domain)

}

// fetchConfig
//
//	@Description: retrieves the models.MailConfig from config/config.json.
//	@return models.MailConfig
func fetchConfig() models.MailConfig {
	var config models.MailConfig

	data, err := os.ReadFile(configFile)

	if len(data) == 0 {
		Logger.Error(GetCurrentFuncName() + " No JSON config data found!")
		return config
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		Logger.Error(GetCurrentFuncName()+" JSON MarshalIndent error!", slog.Any("output", err))
		return config
	}

	return config
}

// SendMail
//
//	@Description: sends a mail to models.TempUser to create his account.
//	@param temp
//	@param status
func SendMail(temp *models.TempUser, status string) {
	// Fetching mail configuration
	config := fetchConfig()

	// Recipient information
	recipientMail := []string{temp.User.Email}

	// Generating confirmation Id
	temp.ConfirmID = generateConfirmationID()

	var subject, templateName string

	switch status {
	case "creation":
		subject = "Email verification"
		templateName = "creation-mail"
	case "lost":
		subject = "Set a new password"
		templateName = "new-password-mail"
	}

	// Setting the headers
	header := make(map[string]string)
	header["From"] = "MangaThorg" + "<" + config.Email + ">"
	header["To"] = temp.User.Email
	header["Subject"] = subject
	header["Message-ID"] = generateMessageID(config.Hostname)
	header["Content-Type"] = "text/html; charset=UTF-8"

	t, err := template.ParseFiles(Path + "templates/" + templateName + ".gohtml")
	if err != nil {
		log.Fatal(err)
	}

	// Create a buffer to hold the formatted message
	var body bytes.Buffer

	// Execute the mail's template with data
	err = t.Execute(&body, struct {
		Username  string
		ConfirmID string
	}{
		Username:  temp.User.Username,
		ConfirmID: temp.ConfirmID,
	})
	if err != nil {
		log.Fatal(err)
	}

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	// Setting the authentication
	auth := smtp.PlainAuth(
		"",
		config.Email,
		config.Auth,
		config.Hostname,
	)

	// Sending the mail using TLS
	err = sendMailTLS(
		fmt.Sprintf("%s:%d", config.Hostname, config.Port),
		auth,
		config.Email,
		recipientMail,
		[]byte(message),
	)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Mail sent successfully!")
	}
}

// dial
//
//	@Description: returns a smtp client.
//	@param addr
//	@return *smtp.Client
//	@return error
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing ApiErr:", err)
		return nil, err
	}
	// Explode Host Port String
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// sendMailTLS: refer to net/smtp func SendMail().
//
// When using net.Dial to connect to the tls (SSL) port, smtp. NewClient() will be stuck and will not prompt err
// When len (to)>1, to [1] starts to prompt that it is secret delivery.
func sendMailTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	// Create smtp client
	c, err := dial(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()

	// Checking authentication
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("ApiErr during AUTH", err)
				return err
			}
		}
	}

	// Setting recipient
	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	// Retrieving the Writer to set the message headers and body
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Writing `msg` in the Writer
	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	// Closing the Writer
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
