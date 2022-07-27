package controller

import (
	"fmt"
	"gitlab.hho-inc.com/devops/flowctl-go/models"
	"gitlab.hho-inc.com/devops/flowctl-go/utils"
	"time"
)

func DumpDeployInfo(env, commit, datetime, git, branch string) {
	config := utils.LoadYaml()
	app := config.GetString("app")
	proj := config.GetString("proj")
	fmt.Printf(`
=================================================================
# 应用：%s
# 环境：%s
# 项目: %s
# 代码版本: %s
# git: %s
# branch: %s
# cpu: 1
# 内存：1G
# 访问地址: %s-%s.hho-inc.com
# 任何问题请联系ops解决，感谢支持！
================================================================`, app, env, proj, commit, git, branch, app, env)

	history := &models.CDHistory{
		App:      app,
		Env:      env,
		CommitID: commit,
		Proj:     proj,
		GitUrl:   git,
		Branch:   branch,
		ImageTag: commit + "-" + datetime,
		ImageUrl: fmt.Sprintf("reg.hho-inc.com/%s-%s/%s:%s",
			proj, env, app, commit+"-"+datetime),
		DeployTime: time.Now(),
	}

	status := &models.CDStatus{
		App:      app,
		Env:      env,
		CommitID: commit,
		Proj:     proj,
		GitUrl:   git,
		Branch:   branch,
		ImageTag: commit + "-" + datetime,
		ImageUrl: fmt.Sprintf("reg.hho-inc.com/%s-%s/%s:%s",
			proj, env, app, commit+"-"+datetime),
	}

	cli := utils.NewDBClient()
	cli.DBInsertHistory(history)
	cli.DBInsertOrUpdateStatus(status)
}
