package modules

import (
	"math/rand"
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func SmtpSender(to_mail string, register bool) string { //register: true인 경우에만 verify number을 제대로 반환함
	// Set up authentication information.
	godotenv.Load()
	smtppassword := os.Getenv("SMTPPW")
	auth := smtp.PlainAuth("", "disjob", smtppassword, "smtp.naver.com")
	var mail_subject, mail_content, verify_num string
	if register {
		mail_subject = "Silver Oasis 회원가입 인증번호"
		verify_num = strconv.Itoa(int(rand.Intn(100000)))
		mail_content = "Silver Oasis 회원가입 인증번호: " + verify_num + "\r\n" +
			"https://pi.imdhson.com/r/" + to_mail + "/" + verify_num + ""
		//http://.com/login/auth/email/mail@imdhson.com/123123321 이런식으로 가게됨
	} else {
		mail_subject = ""
		mail_content = ""
	}
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{to_mail} //슬라이스타입인 이유는 여러명에게 동시에 보낼 수도 있음
	msg := []byte("To: " + to_mail + "\r\n" +
		"From: disjob@naver.com\r\n" + //네이버 smtp 사용시 From: 을 반드시 추가하여야함!
		"Subject: " + mail_subject + "\r\n" +
		"\r\n" +
		mail_content + "\r\n")
	err := smtp.SendMail("smtp.naver.com:587", auth, "disjob@naver.com", to, msg)
	ErrOK(err)
	return verify_num
}
