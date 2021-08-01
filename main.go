package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

func stopProc(done chan<- error) {
	stopDir := "/Users/gozu/projects/jetty-distribution-9.4.43.v20210629/demo-base"
	stopCmd := "java"
	stopArgs := strings.Fields("-jar ../start.jar STOP.PORT=28282 STOP.KEY=secret --stop")
	//stopArgs := strings.Fields("-jar ../start.jar STOP.PORT=28282 STOP.KEY=secret jetty.http.port=8081 jetty.ssl.port=8444")

	// process stop
	cmd := exec.Command(stopCmd, stopArgs...)
	cmd.Dir = stopDir
	// cmd.Env = startEnv
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- stderr ---")
	scanner2 := bufio.NewScanner(stderr)
	for scanner2.Scan() {
		fmt.Println(scanner2.Text())
	}

	fmt.Println("--- stdout ---")
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	done <- nil
	close(done)
}

func main() {
	// Ctrl+Cを受け取る
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	done := make(chan error, 1)
	go stopProc(done)

	select {
	case <-quit:
		fmt.Println("interrup signal accepted.")
	case err := <-done:
		fmt.Println("exit.", err)
	}
}
