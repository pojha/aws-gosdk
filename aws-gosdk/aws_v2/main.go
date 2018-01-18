package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-martini/martini"
)

func main() {
	fmt.Printf("Hello, I am from 'ip.go' file.")
	m := martini.Classic()
	m.Get("/", func() string {
		nt := getUrls()
		nt += "\n"
		nt += getEc2()
		return nt
	})

	m.Get("/all", func() string {
		nt := getAllEc2()
		return nt
	})

	m.Run()
}

func getEc2() string {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command("bash", "ip.sh")
	stdoutIn, _ := cmd.StdoutPipe()
	var errStdout error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr := string(stdoutBuf.Bytes())
	return outStr
}

func getAllEc2() string {
	//fmt.Printf("All.");
	cmd := exec.Command("bash", "ip_all.sh")
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//  fmt.Printf("Output \n",out.String());
	return out.String()
}

func getUrls() string {
	cmd := exec.Command("cat", "urls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}
