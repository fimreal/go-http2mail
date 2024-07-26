# 邮件发送服务

该项目是一个基于 Golang 和 Gin 框架实现的邮件发送服务。用户可以通过 HTTP POST 请求向指定邮箱发送电子邮件，支持多种邮件服务商的 SMTP 配置。

## 特性

目前仅支持 465 和 587 端口的 SMTP 服务器。

- 支持自定邮件服务商，同时内置如 Gmail、Yahoo、Outlook 等主流邮件服务商
- 可配置的 SMTP 服务器和端口
- 随机生成 API 密钥
- 支持 JSON 格式的请求体
- 错误处理与日志记录

## 安装与运行

dockerhub 仓库：https://hub.docker.com/r/epurs/http2mail

### 前提条件

- 确保已安装 [Go](https://golang.org/doc/install) 环境（版本要求：1.16 及以上）。

### 克隆项目

```bash
git clone https://github.com/fimreal/go-http2mail.git
cd go-http2mail
```

### 打包运行服务

在项目目录中运行以下命令安装依赖:

```bash
make
./http2mail [your_api_key]
```

### 直接运行服务

运行以下命令启动邮件发送服务：

```bash
go run main.go [your_api_key]
```

如果没有提供 API 密钥，将会生成一个随机的 API 密钥并打印到控制台。

### Docker 本地 build

```bash
docker build --pull --tag http2mail:latest -f Dockerfile https://github.com/fimreal/go-http2mail.git
docker run -d --name http2mail -P http2mail:latest [your_api_key]
```

### 设置调试模式

可以通过设置环境变量 `DEBUG` 来启用调试模式：

```bash
DEBUG=1 go run main.go [your_api_key]
```

## API 接口

### 请求方式

- **POST** `/`

### 请求体

请求体需为 JSON 格式，示例：

```json
{
  "api_key": "你的API密钥",
  "email": "发件人邮箱",
  "password": "发件人邮箱密码",
  "sender_name": "发件人名称（可选）",
  "recipient": "收件人邮箱",
  "subject": "邮件主题",
  "message": "邮件内容",
  "smtp_server": "自定义SMTP服务器（可选）",
  "smtp_port": "自定义SMTP端口（可选）"
}
```

### 响应

- 成功响应：

```json
{
  "message": "Email sent successfully"
}
```

- 错误响应：

```json
{
  "error": "错误信息"
}
```

### curl

```bash
curl -X POST http://localhost:8080/ \
-F "api_key=your_api_key_here" \
-F "email=aaa@gmail.com" \
-F "password=your_password" \
-F "sender_name=" \  # 可省略
-F "recipient=recipient@example.com" \
-F "subject=Test Email" \
-F "message=This is a test email." \
-F "smtp_server=smtp.example.com" \  # 可选
-F "smtp_port=465"                    # 可选
```

或者使用 json

```bash
curl -X POST http://localhost:8080/send-email \
-H "Content-Type: application/json" \
-d '{
    "api_key": "your_api_key_here",
    "email": "aaa@gmail.com",
    "password": "your_password",
    "sender_name": "", # 可省略
    "recipient": "recipient@example.com",
    "subject": "Test Email",
    "message": "This is a test email.",
    "smtp_server": "smtp.example.com", # 可选
    "smtp_port": "465"                 # 可选
}'
```
