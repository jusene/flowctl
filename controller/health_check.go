package controller

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl-go/utils"
	"net/http"
	"os"
	"strings"
	"time"
)

type HealthCheck struct {
	client *http.Client
	config *viper.Viper
	env string
}

func NewHealthCheck(env string) *HealthCheck {
	config := utils.LoadYaml()
	client := &http.Client{}

	return &HealthCheck{
		client: client,
		config: config,
		env: env,
	}
}

func (c *HealthCheck) HttpCheck(url, path string) error {
	req, err := http.NewRequest("GET", strings.Join([]string{url, path}, ""), nil)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(2)
	}

	resp, err := c.client.Do(req)
	if resp.StatusCode == 200 {
		fmt.Println("--------- > 应用检查成功，服务正常启动...")
		return nil
	} else {
		fmt.Fprint(os.Stderr, "--------- > 应用检查失败，正常重试...")
		return errors.New("应用检查失败")
	}

}

func (c *HealthCheck) RunCheck() {
	if c.config.GetString("runningtime") != "java8" && c.config.GetString("runningtime") != "java11" {
		fmt.Println("非java应用，跳过检查...")
		return
	}
	fmt.Println("让程序跑一会...")
	time.Sleep(30 * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	count := 1
	for {
		<- ticker.C
		fmt.Printf("---------- > 应用第%d次检查: http://%s\n", count, c.config.GetString("app")+"-"+c.env+".hho-inc.com/health/check")
		err := c.HttpCheck(fmt.Sprintf("http://%s", c.config.GetString("app")+"-"+c.env+".hho-inc.com"), "/health/check")
		if err == nil {
			break
		}
		count++
		if count >= 30 {
			fmt.Println("---------- > 应用检查超过30次，中断检查，请k8s dashboard查看错误日志...")
			os.Exit(2)
		}
	}
}