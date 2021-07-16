package shell

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

type Script string

func (s Script) Run() (errorCode int) {
	args := append(strings.Split(Shell, " "), string(s))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	sigint := make(chan os.Signal)
	closeChan := make(chan struct{})
	go func() {
		select {
		case <-closeChan:
			return
		case <-sigint:
			_ = cmd.Process.Signal(os.Interrupt)
		}
	}()
	signal.Notify(sigint, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	defer close(closeChan)
	err := cmd.Run()
	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return -1
	}
	return exitErr.ExitCode()
}

func (s Script) Eval() (output string) {
	args := append(strings.Split(Shell, " "), string(s))
	cmd := exec.Command(args[0], args[1:]...)
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return ""
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	sigint := make(chan os.Signal)
	closeChan := make(chan struct{})
	go func() {
		select {
		case <-closeChan:
			return
		case <-sigint:
			_ = cmd.Process.Signal(os.Interrupt)
		}
	}()
	signal.Notify(sigint, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	defer close(closeChan)
	err = cmd.Start()
	if err != nil {
		return ""
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	_ = cmd.Wait()
	return string(bytes)
}
