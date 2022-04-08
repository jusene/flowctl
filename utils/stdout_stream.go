package utils

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os/exec"
)

func CmdStreamOut(cmd *exec.Cmd) {
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
	}

	for {
		line, err := readerErr.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Print("Error: ", line)
	}
	cmd.Wait()
}
