package main

import (
	//"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"io/ioutil"
	"sync"

	"github.com/go-martini/martini"
)

var wg sync.WaitGroup

func main() {
	fmt.Printf("Hello, I am from 'ip.go' file.")
	m := martini.Classic()
	m.Get("/ip", func() string {
		fileName := "epen_ip"
		//nt := readFromFile("urls")
		nt += "\n"
		nt += readFromFile(fileName)
		getEc2(fileName)
		//getFile(fileName)
		return nt
	})

	m.Get("/", func() string {
		nt := readFromFile("urls")
		nt += "\n"
		nt += readFromFile("epen_ip")
		return nt
	})

	m.Get("/all", func() string {
		nt := getAllEc2()
		return nt
	})

	m.Get("/file", func() string {
		fileName := "epen_ip"
		getEc2(fileName)
		nt := readFromFile(fileName)
		return nt
	})

	m.Run()
}

func getFile(fileName string) {
	go getEc2(fileName)
}

func getEc2(fileName string) {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command("bash", "ip.sh")
	stdoutIn, _ := cmd.StdoutPipe()
	var errStdout error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}

	err = ioutil.WriteFile(fileName, stdoutBuf.Bytes(), 0644)
	//outStr := string(stdoutBuf.Bytes())
	//return outStr
}

func readFromFile(fileName string) string {
	cmd := exec.Command("cat", fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
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
