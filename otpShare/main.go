package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gopkg.in/yaml.v2"
)

var mu sync.RWMutex

// OTPConfig 定义OTP配置项结构
type OTPConfig struct {
	Secret string `yaml:"Secret"`
	Name   string `yaml:"Name"`
}

// 从配置文件中加载OTP配置
func loadOTPConfigs(filename string) ([]OTPConfig, error) {
	// 读取配置文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析配置文件内容
	var configs []OTPConfig
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

// 计算当前时间的OTP码
func generateOTP(secret string) (string, error) {
	otpCode, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", err
	}
	return otpCode, nil
}

func main() {
	// 设置Gin模式为发布模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	router := gin.Default()

	var otpConfigs []OTPConfig

	go func() {
		for {
			mu.Lock()
			// 加载OTP配置
			otpConfig, err := loadOTPConfigs("etc/conf.yaml")
			if err != nil {
				fmt.Println("Error refresh loading OTP configs:", err)
			}
			otpConfigs = otpConfig
			mu.Unlock()
			time.Sleep(60 * time.Second)
		}
	}()

	// 定义路由处理程序
	router.GET("/", func(c *gin.Context) {
		c.String(200, "此服务由淘宝店铺 <Copilot商业订阅> 提供服务, 请勿滥用\n")
		c.String(200, "OTP码30S 刷新一次 若GitHub登陆显示过期请刷新网页\n")
		mu.RLock()
		// 遍历每个OTP配置项，并生成当前OTP码
		for _, config := range otpConfigs {
			otpCode, err := generateOTP(config.Secret)
			if err != nil {
				fmt.Println("Error generating OTP:", err)
				c.String(500, "Internal Server Error")
				return
			}
			// 将密钥、账号名称和当前OTP码返回给客户端
			c.String(200, "GitHub 邮箱: %s --- OTP密钥: %s\n", config.Name, otpCode)
		}
		mu.RUnlock()
	})

	// 启动HTTP服务器
	err := router.Run(":20024")
	if err != nil {
		panic(err)
	}
}
