package controller

import (
	"fmt"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HHOBuildCode struct {
	config    *viper.Viper
	workSpace string
}

func NewHHOBuildCode() *HHOBuildCode {
	config := utils.LoadYaml()
	currentPath, _ := filepath.Abs(".")
	return &HHOBuildCode{
		config:    config,
		workSpace: currentPath,
	}
}

func (h *HHOBuildCode) Build() {
	fmt.Println("-------------> build code", h.config.Get("app").(string))
	switch strings.ToLower(h.config.GetString("runningtime")) {
	case "java8", "java11":
		fmt.Println("****** 检测应用是java应用 ******")
		h.buildJava()

	case "node":
		fmt.Println("****** 检测应用是node应用 ******")
		h.buildNode()

	case "static":
		fmt.Println("****** 检测应用是static静态页面 ******")
		h.buildStatic()
	default:
		fmt.Println("unknown runningtime, please check app.yaml, runningtime must java8, java11, node, static")
		os.Exit(2)
	}
}

func (h *HHOBuildCode) buildJava() {
	os.Chdir(h.workSpace)
	cmd := exec.Command("mvn", "-B", "clean", "package", "-Dmaven.test.skip=true")
	utils.CmdStreamOut(cmd)
}

func (h *HHOBuildCode) buildNode() {
	os.Chdir(h.workSpace)
	cmd := exec.Command("npm", "install")
	utils.CmdStreamOut(cmd)
}

func (h *HHOBuildCode) buildStatic() {
	os.Chdir(h.workSpace)
	cmd := exec.Command("npm", "install")
	utils.CmdStreamOut(cmd)

	cmd2 := exec.Command("npm", "run", "build")
	utils.CmdStreamOut(cmd2)
}
