package controller

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl/models"
	"gitlab.hho-inc.com/devops/flowctl/utils"
	"os"
	"os/exec"
	"path/filepath"
)

type HHOPublishImage struct {
	config    *viper.Viper
	workSpace string
	env       string
	id        string
	time      string
}

func NewPublishImage(env, id, time string) *HHOPublishImage {
	config := utils.LoadYaml()
	currentPath, _ := filepath.Abs(".")
	return &HHOPublishImage{
		config:    config,
		workSpace: currentPath,
		env:       env,
		id:        id,
		time:      time,
	}
}

func (h *HHOPublishImage) Publish() {
	os.MkdirAll("/tmp/deploy", 0755)
	deployment, err := os.OpenFile(fmt.Sprintf("/tmp/deployment-%s.yaml",
		h.config.GetString("app")), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer deployment.Close()

	appInfo := &models.AppInfo{
		APP:  h.config.GetString("app"),
		PROJ: h.config.GetString("proj"),
		ENV:  h.env,
		TIME: h.time,
		ID:   h.id,
	}

	c := utils.NewConsul()
	c.Render2file("/devops/cicd/build/deployment.yaml", deployment, appInfo)

	cmd := exec.Command("/usr/local/bin/kubectl", "apply", "-f",
		fmt.Sprintf("/tmp/deployment-%s.yaml", h.config.GetString("app")))
	utils.CmdStreamOut(cmd)
}
