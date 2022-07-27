package controller

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl-go/utils"
	"os"
	"time"
)

type HealthCheck struct {
	env string
	conf *viper.Viper
}

func NewHealthCheck(env string) *HealthCheck {
	conf := utils.LoadYaml()
	return &HealthCheck{
		env: env,
		conf: conf,
	}
}

func (c *HealthCheck) Check() (replicas, availableReplicas int32) {
	cli := utils.NewKubeClient(c.env)
	replicas, availableReplicas = cli.ListDeployment()
	return
}

func (c *HealthCheck) RunCheck() {
	fmt.Println("让程序跑一会...")
	time.Sleep(30 * time.Second)
	ticker := time.NewTicker(5 * time.Second)
	count := 1
	for {
		<- ticker.C
		fmt.Printf("---------- > 应用第%d次检查:\n", count)
		replicas, availableReplicas := c.Check()
		if replicas != availableReplicas {
			fmt.Printf("Deployment: %s, Replicas: %d, AvaileRelicas: %d\n", c.conf.GetString("app"), replicas, availableReplicas)
		} else {
			fmt.Println("服务启动成功")
			break
		}
		count++
		if count >= 30 {
			fmt.Println("---------- > 应用检查超过30次，中断检查，请k8s dashboard https://k8s.hho-inc.com 查看错误日志...")
			os.Exit(2)
		}
	}
}