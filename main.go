package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

func stopProc(done chan<- error) {
	stopDir := "/Users/gozu/projects/jetty-distribution-9.4.43.v20210629/demo-base"
	stopCmd := "java"
	//stopArgs := strings.Fields("-jar ../start.jar STOP.PORT=28282 STOP.KEY=secret --stop")
	stopArgs := strings.Fields("-jar ../start.jar STOP.PORT=28282 STOP.KEY=secret jetty.http.port=8081 jetty.ssl.port=8444")

	// process stop
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, stopCmd, stopArgs...)
	cmd.Dir = stopDir
	// cmd.Env = startEnv
	log.Println("--- stop cmd start ---")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("--- stdout & stderr ---")
	log.Printf("%s\n", stdoutStderr)
	log.Println("stop process")

	done <- nil
	close(done)
}

func stopProcByPid(done chan<- error, pid int) {
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
	}
	err = p.Kill()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("stop process")

	done <- nil
	close(done)
}

func main() {
	// Ctrl+Cを受け取る
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	done := make(chan error, 1)
	//go stopProc(done)
	pid := 22842
	go stopProcByPid(done, pid)

	select {
	case <-quit:
		fmt.Println("interrup signal accepted.")
	case err := <-done:
		fmt.Println("exit.", err)
	}
}
