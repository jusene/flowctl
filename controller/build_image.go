package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl-go/models"
	"gitlab.hho-inc.com/devops/flowctl-go/utils"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type HHOBuildImage struct {
	config    *viper.Viper
	workSpace string
	env       string
	id        string
	time      string
}

func NewHHOBuildImage(env, id, time string) *HHOBuildImage {
	config := utils.LoadYaml()
	currentPath, _ := filepath.Abs(".")
	return &HHOBuildImage{
		config:    config,
		workSpace: currentPath,
		env:       env,
		id:        id,
		time:      time,
	}
}

func (h *HHOBuildImage) Build() {
	app := h.config.GetString("app")
	switch strings.ToLower(h.config.GetString("runningtime")) {
	case "java8", "java11":
		h.preJavaBuild(app)
	case "node":
		h.preNodeBuild(app)
	case "static":
		h.preStaticBuild(app)
	case "golang":
		h.preGolangBuild(app)
	}
	h.buildImage(app)
	h.pushImage(app)
}

func (h *HHOBuildImage) preJavaBuild(app string) {
	os.Chdir(h.workSpace)
	workDir := strings.Join([]string{h.workSpace, "docker"}, "/")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建%s目录失败", strings.Join([]string{h.workSpace, "docker"}, "/"))
	}

	var jarPath = strings.Join([]string{h.workSpace, app + ".starter", "target", app + ".jar"}, "/")
	if h.config.GetString("jarpath") != "" {
		jarPath = strings.Join([]string{h.workSpace, strings.Trim(h.config.GetString("jarpath"), "/")}, "/")
	}

	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s jar包不存在，请检查生成jar包规则路径...", strings.Join([]string{h.workSpace, app + ".starter", "target", app + ".jar"}, "/"))
		os.Exit(2)
	}

	destFile, _ := os.Create(strings.Join([]string{workDir, app + ".jar"}, "/"))
	defer destFile.Close()
	srcFile, _ := os.Open(jarPath)
	defer srcFile.Close()
	io.Copy(destFile, srcFile)

	destFile, _ = os.Create(strings.Join([]string{workDir, "app.yaml"}, "/"))
	srcFile, _ = os.Open(strings.Join([]string{h.workSpace, "app.yaml"}, "/"))
	io.Copy(destFile, srcFile)

	dockerfile, err := os.OpenFile(workDir+"/Dockerfile", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer dockerfile.Close()

	c := utils.NewConsul()
	appInfo := &models.AppInfo{
		APP:      app,
		RUNNTIME: h.config.GetString("runningtime"),
		ENV:      h.env,
		DEBPACK:   h.config.GetString("debpack"),
	}

	c.Render2file("/devops/cicd/build/dockerfile", dockerfile, appInfo)
}

func (h *HHOBuildImage) preNodeBuild(app string) {
	os.Chdir(h.workSpace)

	dockerfile, err := os.OpenFile(h.workSpace+"/Dockerfile", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer dockerfile.Close()

	c := utils.NewConsul()
	appInfo := &models.AppInfo{
		APP:      app,
		RUNNTIME: h.config.GetString("runningtime"),
		DEBPACK:   h.config.GetString("debpack"),
		ENV:      h.env,
	}
	c.Render2file("/devops/cicd/build/dockerfile", dockerfile, appInfo)
}

func (h *HHOBuildImage) preGolangBuild(app string) {
	os.Chdir(h.workSpace)
	dockerfile, err := os.OpenFile(h.workSpace+"/Dockerfile", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer dockerfile.Close()

	c := utils.NewConsul()
	appInfo := &models.AppInfo{
		APP:      app,
		RUNNTIME: h.config.GetString("runningtime"),
		DEBPACK:   h.config.GetString("debpack"),
		ENV:      h.env,
	}
	c.Render2file("/devops/cicd/build/dockerfile", dockerfile, appInfo)
}

func (h *HHOBuildImage) preStaticBuild(app string) {
	os.Chdir(h.workSpace)
	workDir := strings.Join([]string{h.workSpace, "docker"}, "/")
	if err := os.MkdirAll(workDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建%s目录失败", strings.Join([]string{h.workSpace, "docker"}, "/"))
	}

	filepath.Walk("build", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			os.MkdirAll(strings.Join([]string{workDir, "build", path}, "/"), 0755)
		} else {
			destFile, _ := os.Create(strings.Join([]string{workDir, "build", path}, "/"))
			defer destFile.Close()
			srcFile, _ := os.Open(path)
			defer srcFile.Close()
			io.Copy(destFile, srcFile)
		}
		return nil
	})

	dockerfile, err := os.OpenFile(workDir+"/Dockerfile", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer dockerfile.Close()

	c := utils.NewConsul()
	appInfo := &models.AppInfo{
		APP:      app,
		RUNNTIME: h.config.GetString("runningtime"),
		DEBPACK:   h.config.GetString("debpack"),
	}
	c.Render2file("/devops/cicd/build/dockerfile", dockerfile, appInfo)
}

func (h *HHOBuildImage) buildImage(app string) {
	fmt.Println("-------------> Build Image ", app)
	switch strings.ToLower(h.config.GetString("runningtime")) {
	case "java8", "java11", "node":
		os.Chdir("./docker")
	case "static":
	case "golang":
		
	}
	utils.CmdStreamWithErr("docker -H tcp://127.0.0.1:2376 build -t "+
		fmt.Sprintf("reg.hho-inc.com/%s-%s/%s:%s", h.config.GetString("proj"), h.env, app, h.id+"-"+h.time) + " .")
}

func (h *HHOBuildImage) pushImage(app string) {
	const DockerHOST = "tcp://127.0.0.1:2376"
	os.Setenv("DOCKER_HOST", DockerHOST)
	fmt.Println("-------------> Push Image ", app)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	cobra.CheckErr(err)

	authConfig := types.AuthConfig{
		Username: "admin",
		Password: "hhoroot@2022",
	}
	encodedJSON, err := json.Marshal(authConfig)
	cobra.CheckErr(err)

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, err := cli.ImagePush(context.Background(), fmt.Sprintf("reg.hho-inc.com/%s-%s/%s:%s",
		h.config.GetString("proj"), h.env, app, h.id+"-"+h.time),
		types.ImagePushOptions{RegistryAuth: authStr})
	cobra.CheckErr(err)
	defer out.Close()
	io.Copy(os.Stdout, out)
}
