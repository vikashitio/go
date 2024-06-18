package function

import (
	"net/smtp"
	"os"
)

type MyStatus struct {
	Status string
}

func GetStatus(Status int) (MyStatus, error) {
	myStatus := ""

	if Status == 1 {
		myStatus = "Active"

	} else if Status == 2 {
		myStatus = "Inactive"
	} else if Status == 3 {
		myStatus = "New"
	} else {
		myStatus = "Deleted"

	}

	//fmt.Println(myStatus)
	var data = MyStatus{
		Status: myStatus,
	}
	return data, nil
}
func SendEmail(subject, HTMLbody string) error {
	// sender data
	var Email = "vikashg@itio.in"

	// smtp - Details
	var fromEmail = os.Getenv("SMTPusername")
	var SMTPpassword = os.Getenv("SMTPpassword")
	var EntityName = os.Getenv("SMTPsendername")
	host := os.Getenv("SMTPhost")
	port := os.Getenv("SMTPport")
	address := host + ":" + port

	to := []string{Email}
	// Set up authentication information.
	auth := smtp.PlainAuth("", fromEmail, SMTPpassword, host)
	msg := []byte(
		"From: " + EntityName + ": <" + fromEmail + ">\r\n" +
			"To: " + Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME: MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			HTMLbody)
	err := smtp.SendMail(address, auth, fromEmail, to, msg)
	if err != nil {
		return err
	}
	//fmt.Println("Check for sent email!")
	return nil
}
