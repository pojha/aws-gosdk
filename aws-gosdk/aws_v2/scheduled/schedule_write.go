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
	"time"
	"net/http"
)

var wg sync.WaitGroup

//GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o schedule.linux main.go
func main() {
	
	go doWriteToFile(20*time.Minute)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,  getEpenEc2())
	})
	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,  getAllEc2())
	})

	http.HandleFunc("/alls", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,  readFromFile("all_ip"))
	})

	http.HandleFunc("/epen", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,  readFromFile("epen_ip"))
	})
	
	fmt.Printf("Hello, I am from 'ip.go' file. 3000")
	http.ListenAndServe(":3000", nil)
	
}

func getEpenEc2() string {
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
	fmt.Printf("All.");
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

func readFromFile(fileName string) string {
	cmd := exec.Command("cat", fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}




func writeToFile(scriptName string, fileName string) string {
	var stdoutBuf bytes.Buffer
	cmd := exec.Command("bash", scriptName)
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
	outStr := string(stdoutBuf.Bytes())

	err = ioutil.WriteFile(fileName, stdoutBuf.Bytes(), 0644)
	return outStr
}

func doWriteToFile(d time.Duration) {
	//fmt.Println("doWriteToFile Tick at", d)
	ticker := time.NewTicker(time.Hour * 2)
	 defer ticker.Stop()
    go func() {
        for t := range ticker.C {
		   writeToFile("ip.sh", "epen_ip")
		   writeToFile("ip_all.sh", "all_ip")
		     fmt.Println("Tick at", t)
        }
    }()
	time.Sleep(time.Minute * 5)	
}