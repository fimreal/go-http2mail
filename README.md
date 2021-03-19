
#### 用途
接收写有邮件内容的 POST 请求，用配置的邮箱发送邮件。

例如：

```bash
curl localhost:5000/mailme -X POST -d 'mailto=abc@abc.com' -d 'subject=邮件标题' -d "body=邮件内容"
```



#### Docker 启动说明

**配置办法一**：

可以将配置文件 go-http2mail.yaml 映射到 "/" ，程序启动时会自动读取。

**配置办法二**：

使用环境变量。

环境变量可选配置：

```
ENV PORT=5000
USERNAME=postuser@mail.domain
PASSWORD=postuserpassword
SMTPSERVER=smtp.mail.com
SMTPSERVERPORT=25
```

**配置办法三**：

启动传入参数

    Usage of go-http2mail:
      --password string         设置发件人的账号密码 (default "xxxx")
      --port string             配置启动监听端口 (default "5000")
      --smtpserver string       设置发件人使用的邮件服务器 (default "smtp.mail.com")
      --smtpserverport string   设置发件人使用的邮件服务器端口 (default "25")
      --username string         设置发件人的账号名称 (default "postuser@mail.com")
三种办法配置优先级由下到上。

