package controller

import (
	"fmt"
	"gitlab.hho-inc.com/devops/flowctl-go/utils"
	"strings"
)

func RollbackImage(env string) {
	config := utils.LoadYaml()
	app := config.GetString("app")
	proj := config.GetString("proj")

	switch strings.ToLower(config.GetString("runningtime")) {

	case "static":
		fmt.Println("静态资源在oss上，请联系管理员处理")
	default:
		rollback(app, proj, env)
	}
}


func rollback(app, proj, env string) {
	utils.CmdStreamOut("kubectl rollout undo deployment " + app +" -n " + fmt.Sprintf("%s-%s", proj, env))
}

