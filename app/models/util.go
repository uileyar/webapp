package models

import (
	"crypto/tls"
	"math"

	"github.com/golang/glog"
	"github.com/satori/go.uuid"

	"github.com/Unknwon/goconfig"
	"gopkg.in/gomail.v2"
)

func CreateGUID() string {
	u1 := uuid.NewV4()
	glog.Infof("UUID: %s\n", u1)
	return u1.String()
}

func Round(f float32, n int) float32 {
	pow10_n := math.Pow10(n)
	return float32(math.Trunc((float64(f)+0.5/pow10_n)*pow10_n) / pow10_n)
}

func SendMail(subject string, body string) {
	var session string = "mail"

	c, err := goconfig.LoadConfigFile("webapp.ini")
	if err != nil {
		panic(err)
	}

	from, _ := c.GetValue(session, "from")
	tolist := c.MustValueArray(session, "to", ",")
	host, _ := c.GetValue(session, "host")
	port, _ := c.Int(session, "port")
	username, _ := c.GetValue(session, "username")
	password, _ := c.GetValue(session, "password")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", tolist...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

type SendConfirmationEmail struct {
	Subject string
	Body    string
}

func (e SendConfirmationEmail) Run() {
	// 查询数据库
	// 发送电子邮件
	SendMail(e.Subject, e.Body)
}
