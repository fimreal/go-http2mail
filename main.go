package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// EmailRequest 定义请求体结构
type EmailRequest struct {
	APIKey     string `json:"api_key"`               // API 密钥
	Username   string `json:"email"`                 // 发件人邮箱
	Password   string `json:"password"`              // 发件人邮箱密码
	SenderName string `json:"sender_name"`           // 发件人名称
	Recipient  string `json:"recipient"`             // 收件人邮箱
	Subject    string `json:"subject"`               // 邮件主题
	Message    string `json:"message"`               // 邮件内容
	SMTPServer string `json:"smtp_server,omitempty"` // SMTP 服务器地址 (可选)
	SMTPPort   string `json:"smtp_port,omitempty"`   // SMTP 端口 (可选)
}

// SMTPConfig 定义邮箱服务商的配置结构
type SMTPConfig struct {
	Server string
	Port   string
}

// Global variable for API key
var apiKey string

// 邮件服务商 SMTP 配置映射
var smtpConfigs = map[string]SMTPConfig{
	"gmail.com":    {"smtp.gmail.com", "587"},
	"yahoo.com":    {"smtp.mail.yahoo.com", "587"},
	"outlook.com":  {"smtp-mail.outlook.com", "587"},
	"hotmail.com":  {"smtp-mail.outlook.com", "587"},
	"163.com":      {"smtp.163.com", "465"},
	"126.com":      {"smtp.126.com", "465"},
	"yeah.com":     {"smtp.yeah.com", "465"},
	"qq.com":       {"smtp.qq.com", "587"},
	"sina.com.cn":  {"smtp.sina.com.cn", "465"},
	"me.com":       {"smtp.mail.me.com", "587"},
	"zoho.com":     {"smtp.zoho.com", "587"},
	"mailgun.org":  {"smtp.mailgun.org", "587"},
	"sendgrid.net": {"smtp.sendgrid.net", "587"},
}

// generateRandomAPIKey 生成随机的 API 密钥
func generateRandomAPIKey() string {
	b := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(b); err != nil {
		fmt.Println("Error generating random API key:", err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

// getSMTPConfig 根据发件人邮箱获取相应的 SMTP 服务器和端口
func getSMTPConfig(email string) (smtpServer, smtpPort string, err error) {
	emailDomain := strings.Split(email, "@")[1]

	if config, exists := smtpConfigs[emailDomain]; exists {
		return config.Server, config.Port, nil
	}

	return "", "", fmt.Errorf("unsupported email domain: %s", emailDomain)
}

// sendEmail 发送电子邮件的函数
func sendEmail(req EmailRequest) error {
	var smtpServer, smtpPort string
	var err error

	// 检查是否提供了自定义的 SMTP 配置
	if req.SMTPServer != "" && req.SMTPPort != "" {
		smtpServer = req.SMTPServer
		smtpPort = req.SMTPPort
	} else {
		smtpServer, smtpPort, err = getSMTPConfig(req.Username)
		if err != nil {
			return err
		}
	}

	var auth smtp.Auth
	addr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	// 根据端口号选择合适的发送方式
	switch smtpPort {
	case "465": // SSL方式
		// 使用 tls.Dial 进行 SSL 连接
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpServer,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		c, err := smtp.NewClient(conn, smtpServer)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

		auth = smtp.PlainAuth("", req.Username, req.Password, smtpServer)

		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
		defer c.Close()

		to := []string{req.Recipient}
		msg := []byte(fmt.Sprintf("From: %s <%s>\nTo: %s\nSubject: %s\n\n%s",
			req.SenderName, req.Username, req.Recipient, req.Subject, req.Message))

		if err = c.Mail(req.Username); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		for _, recipient := range to {
			if err = c.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to set recipient: %w", err)
			}
		}

		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		defer w.Close()

		if _, err = w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

	case "587": // STARTTLS方式
		conn, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}

		if ok, _ := conn.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				InsecureSkipVerify: true,
				ServerName:         smtpServer,
			}
			if err = conn.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("failed to start TLS: %w", err)
			}
		}

		auth = smtp.PlainAuth("", req.Username, req.Password, smtpServer)

		if err = conn.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		to := []string{req.Recipient}
		msg := []byte(fmt.Sprintf("From: %s <%s>\nTo: %s\nSubject: %s\n\n%s",
			req.SenderName, req.Username, req.Recipient, req.Subject, req.Message))

		if err = conn.Mail(req.Username); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		for _, recipient := range to {
			if err = conn.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to set recipient: %w", err)
			}
		}

		w, err := conn.Data()
		if err != nil {
			return fmt.Errorf("failed to get data writer: %w", err)
		}
		defer w.Close()

		if _, err = w.Write(msg); err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

	default:
		return fmt.Errorf("unsupported SMTP port: %s", smtpPort)
	}

	return nil
}

// emailHandler 处理发送邮件的 HTTP 请求
func emailHandler(c *gin.Context) {
	var emailReq EmailRequest
	if err := c.ShouldBindJSON(&emailReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Println("Invalid request body:", err)
		return
	}

	// 验证 API 密钥
	if apiKey != "" && emailReq.APIKey != apiKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		log.Printf("Unauthorized access with API Key: %s\n", emailReq.APIKey)
		return
	}

	if err := sendEmail(emailReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send email: %v", err)})
		log.Printf("Failed to send email to %s: %v\n", emailReq.Recipient, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
	log.Printf("Email sent successfully to %s\n", emailReq.Recipient)
}

func main() {
	// 从命令行参数获取 API 密钥
	if len(os.Args) > 1 {
		apiKey = os.Args[1]
	} else {
		apiKey = generateRandomAPIKey()
	}

	fmt.Printf("Generate Random API Key: %s\n", apiKey)

	if os.Getenv("DEBUG") != "" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.POST("/", emailHandler)

	log.Println("Starting server on :8782")
	log.Fatal(r.Run(":8782"))
}
