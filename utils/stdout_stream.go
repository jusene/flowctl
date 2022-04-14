package utils

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"strings"
)

func CmdStreamOut(cmd *exec.Cmd) {
	errorChan := make(chan string)
	stdout, err := cmd.StdoutPipe()
	cobra.CheckErr(err)
	stderr, err := cmd.StderrPipe()
	cobra.CheckErr(err)
	cmd.Start()

	reader := bufio.NewReader(stdout)
	readerErr := bufio.NewReader(stderr)

	for {
		// 以换行符作为一行结尾
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print(line)
		go func() {
			if strings.Contains(strings.ToLower(line), "error") {
				errorChan <- line
			}
		}()
	}

	for {
		line, err := readerErr.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print("Error: ", line)
		go func() {
			if strings.Contains(strings.ToLower(line), "error") {
				errorChan <- line
			}
		}()
	}
	cmd.Wait()

	if len(errorChan) != 0 {
		fmt.Print("-------------------> 错误详情")
		for msg := range errorChan {
			fmt.Print(msg)
		}
		close(errorChan)
		os.Exit(2)
	}
}
