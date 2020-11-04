package config

// package config

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// 命令行标志定义
var (
	// conf     = flag.String("f", "./go-http2mail", "指定使用的配置文件")
	port           = pflag.String("port", "5000", "配置启动监听端口")
	username       = pflag.String("username", "postuser@mail.com", "设置发件人的账号名称")
	password       = pflag.String("password", "xxxx", "设置发件人的账号密码")
	smtpserver     = pflag.String("smtpserver", "smtp.mail.com", "设置发件人使用的邮件服务器")
	smtpserverport = pflag.String("smtpserverport", "25", "设置发件人使用的邮件服务器端口")
)

// BindEnvFor 绑定环境变量
func BindEnvFor() {
	// 绑定环境变量
	viper.BindEnv("port", "PORT")
	viper.BindEnv("username", "USERNAME")
	viper.BindEnv("password", "PASSWORD")
	viper.BindEnv("smtpserver", "SMTPSERVER")
	viper.BindEnv("smtpserverport", "SMTPSERVERPORT")
}

// init 初始化从执行文件所在目录查找配置文件并加载
func init() {
	// 设置 port 默认值 5000, 前面命令行配置已经有了
	// viper.SetDefault("port", 5000)
	// viper.SetDefault("smtpserverport", 25)

	// 配置文件名字，不包含后缀。 此外可以手动指定配置文件格式： viper.SetConfigType("yaml")
	viper.SetConfigName("go-mail")
	// 添加配置搜索的第一个路径，设置为与二进制文件同目录
	viper.AddConfigPath(".")
	// 判断加载配置文件是否正确
	if err := viper.ReadInConfig(); err != nil {
		// 判断是否是因为找不到文件
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 如果是因为找不到文件，则忽略该错误
			log.Printf("Warning: %v", err)
		} else {
			// 如果是因为文件读取出现错误，则报错退出
			log.Fatalf("Read config file failed: %v", err)
		}
	}

}

// LoadConfigs 加载配置
func LoadConfigs() {
	// 加载环境变量
	BindEnvFor()
	// 解析传入参数
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
