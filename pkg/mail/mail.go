package mail

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/viper"

	"gopkg.in/gomail.v2"
)

// Mailsrv 配置发件人信息
func Mailsrv(mailTo []string, subject string, body string) error {
	mailConn := map[string]string{
		"user": viper.GetString("username"),
		"pass": viper.GetString("password"),
		"host": viper.GetString("smtpserver"),
		"port": viper.GetString("smtpserverport"),
	}

	port, _ := strconv.Atoi(mailConn["port"])

	m := gomail.NewMessage()

	// m.SetHeader("From", mailConn["user"])
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "GoMail Robot"))
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	// 解决 x509: certificate signed by unknown authority 报错问题
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(m)
	return err
}

// SendMail 用来配置邮件收件人、标题、内容, 返回 http.Response 格式
func SendMail(w http.ResponseWriter, r *http.Request) {
	// 判断如果不是 POST 则报错
	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		return
	}

	// 解析url传递的参数，对于POST则解析响应包的主体（request body）
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		log.Println("ParseForm() err: ", err)
		return
	}

	// 收件人
	mailTo := []string{r.FormValue("mailto")}

	// 邮件标题
	subject := r.FormValue("subject")

	// 邮件内容
	body := r.FormValue("body")

	err := Mailsrv(mailTo, subject, body)
	if err != nil {
		// 日志打印失败信息
		log.Println(err, "\n    Details: {\"mailTo\": \"", mailTo, "\", \"subject\": \"", subject, "\"}")
		// 返回失败信息
		fmt.Fprintln(w, "Send fail!")
		return
	}
	// 日志打印成功信息
	log.Println("Send successfully!\n    Details: {\"mailTo\": \"", mailTo, "\", \"subject\": \"", subject, "\"}")
	// 返回成功信息
	fmt.Fprintln(w, "Send successfully!")
}
