package controller

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl/utils"
	"io"
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
	stdout, err := cmd.StdoutPipe()
	cobra.CheckErr(err)
	cmd.Start()

	// 创建一个流来读取管道内的内容，一行一行读
	reader := bufio.NewReader(stdout)

	for {
		// 以换行符作为一行结尾
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}
	cmd.Wait()
}

func (h *HHOBuildCode) buildNode() {
	os.Chdir(h.workSpace)
	cmd := exec.Command("npm", "install")
	stdout, err := cmd.StdoutPipe()
	cobra.CheckErr(err)
	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		// 以换行符作为一行结尾
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
	}
	cmd.Wait()
}

func (h *HHOBuildCode) buildStatic() {
	os.Chdir(h.workSpace)
	cmd := exec.Command("npm", "install")
	stdout, err := cmd.StdoutPipe()
	cobra.CheckErr(err)
	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		// 以换行符作为一行结尾
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
	}
	cmd.Wait()

	cmd2 := exec.Command("npm", "run", "build")
	stdout2, err := cmd2.StdoutPipe()
	cobra.CheckErr(err)
	cmd.Start()

	reader2 := bufio.NewReader(stdout2)

	for {
		// 以换行符作为一行结尾
		line, err := reader2.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
	}
	cmd.Wait()
}
