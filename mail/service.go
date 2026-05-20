package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Mail struct {
	To      []string
	Subject string
	Body    string
}

func (mail *Mail) BuildMessage() []byte {
	message := "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	message += fmt.Sprintf("From: PhD Recruitment Automation System IITK <%s>\r\n", sender)
	message += fmt.Sprintf("Subject: %s | PhD Recruitment Automation System\r\n", mail.Subject)

	if len(mail.To) == 1 {
		message += fmt.Sprintf("To: %s\r\n\r\n", mail.To[0])
	} else {
		message += fmt.Sprintf("To: Undisclosed Recipients <%s>\r\n\r\n", webteam)
	}

	message += strings.Replace(mail.Body, "\n", "<br>", -1)
	message += "<br><br>--<br>PhD Recruitment Automation System<br>"
	message += "Indian Institute of Technology Kanpur<br><br>"
	message += "<small>This is an auto-generated email. Please do not reply.</small>"

	return []byte(message)
}

func batchEmails(to []string, batch int) [][]string {
	var batches [][]string

	for i := 0; i < len(to); i += batch {
		end := i + batch
		if end > len(to) {
			end = len(to)
		}
		batches = append(batches, to[i:end])
	}

	return batches
}

// For IITK mmtp.iitk.ac.in:465
func sendMailTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: host, // mmtp.iitk.ac.in
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(msg)
	if err != nil {
		return err
	}

	return writer.Close()
}

func Service(mailQueue chan Mail) {
	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	logrus.Infof(
	"SMTP config host=%s port=%s user=%s sender=%s pass_len=%d",
	host, port, user, sender, len(pass),
)
	for mail := range mailQueue {
		message := mail.BuildMessage()

		to := append(mail.To, webteam)
		batches := batchEmails(to, batch)

		for _, emailBatch := range batches {
			if err := sendMailTLS(addr, auth, sender, emailBatch, message); err != nil {
				logrus.Errorf("Error sending mail to: %v", emailBatch)
				logrus.Errorf("Error: %v", err)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
