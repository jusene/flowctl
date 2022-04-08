package controller

import (
	"fmt"
	"gitlab.hho-inc.com/devops/flowctl/utils"
)

func DumpDeployInfo(env string) {
	config := utils.LoadYaml()
	app := config.GetString("app")
	proj := config.GetString("proj")
	fmt.Printf(`
=================================================================
# 应用：%s
# 环境：%s
# 项目: %s
# cpu: 0.5
# 内存：512M
# 访问地址: %s-%s.hho-inc.com
# 任何问题请联系ops解决，感谢支持！
================================================================`, app, env, proj, app, env)
}
