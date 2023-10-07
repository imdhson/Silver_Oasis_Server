package modules

import (
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func CriticalAlert(err error) {
	// Set up authentication information.
	godotenv.Load()
	smtppassword := os.Getenv("SMTPPW")
	auth := smtp.PlainAuth("", "disjob", smtppassword, "smtp.naver.com")
	var mail_subject, mail_content string
	mail_subject = "DisJob 서버가 꺼짐!!!"
	mail_content = err.Error() + "의 이유로 서버가 꺼짐."

	to_mail := "hhammer1234@gmail.com"

	to := []string{to_mail}
	msg := []byte("To: " + to_mail + "\r\n" +
		"From: disjob@naver.com\r\n" + //네이버 smtp 사용시 From: 을 반드시 추가하여야함!
		"Subject: " + mail_subject + "\r\n" +
		"\r\n" +
		mail_content + "\r\n")
	err = smtp.SendMail("smtp.naver.com:587", auth, "disjob@naver.com", to, msg)
	ErrOK(err)
}
