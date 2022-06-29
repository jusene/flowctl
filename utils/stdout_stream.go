package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func CmdStreamOut(cmd string)  {
	ctx, cancel := context.WithCancel(context.Background())
	if err := Command(ctx, cmd); err != nil {
		fmt.Print(err.Error())
		cancel()
		os.Exit(2)
	}
}

func CmdStreamWithErr(cmd string) {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan string, 1)

	go func(errChan chan string, cancel context.CancelFunc) {
		for {
			e := <- errChan
				fmt.Sprintln(os.Stderr, "发生错误: ", e)
				cancel()
				os.Exit(2)
		}
	}(errChan, cancel)

	if err := CommandWithErr(ctx, cmd, errChan); err != nil {
		fmt.Print(err.Error())
		cancel()
		os.Exit(2)
	}
}

func Command(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go read(ctx, &wg, stderr)
	go read(ctx, &wg, stdout)

	err = c.Start()
	wg.Wait()
	return err
}

func CommandWithErr(ctx context.Context, cmd string, errChan chan string) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go readWithErr(ctx, &wg, stderr, errChan)
	go readWithErr(ctx, &wg, stdout, errChan)

	err = c.Start()
	wg.Wait()
	return err
}

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <- ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}
}


func readWithErr(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser, errChan chan string) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <- ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
			go func(readString string) {
				if strings.Contains(strings.ToLower(readString), "error") {
					errChan <- readString
				}
			}(readString)
		}
	}
}