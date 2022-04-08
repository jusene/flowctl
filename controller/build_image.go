package controller

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.hho-inc.com/devops/flowctl/models"
	"gitlab.hho-inc.com/devops/flowctl/utils"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HHOBuildImage struct {
	config    *viper.Viper
	workSpace string
	env       string
}

func NewHHOBuildImage(env string) *HHOBuildImage {
	config := utils.LoadYaml()
	currentPath, _ := filepath.Abs(".")
	return &HHOBuildImage{
		config:    config,
		workSpace: currentPath,
		env:       env,
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

	destFile, _ := os.Create(strings.Join([]string{workDir, app + ".jar"}, "/"))
	defer destFile.Close()
	srcFile, _ := os.Open(strings.Join([]string{h.workSpace, app + ".starter", "target", app + ".jar"}, "/"))
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
	}
	c.Render2file("/devops/cicd/build/dockerfile", dockerfile, appInfo)
}

func (h *HHOBuildImage) buildImage(app string) {
	fmt.Println("-------------> Build Image ", app)
	switch strings.ToLower(h.config.GetString("runningtime")) {
	case "java8", "java11", "node":
		os.Chdir("./docker")
	case "static":
	}
	cmd := exec.Command("docker", "-H", "tcp://127.0.0.1:2376", "build", "-t",
		fmt.Sprintf("reg.hho-inc.com/%s-%s/%s:%s", h.config.GetString("proj"), h.env, app, "latest"), ".")
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
		h.config.GetString("proj"), h.env, app, "latest"),
		types.ImagePushOptions{RegistryAuth: authStr})
	cobra.CheckErr(err)
	defer out.Close()
	io.Copy(os.Stdout, out)
}
