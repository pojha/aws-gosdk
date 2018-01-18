package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/go-martini/martini"
)

func main() {
	fmt.Printf("main_shell, I am from 'ip.go' file.")
	m := martini.Classic()

	m.Get("/", func() string {
		nt := getEc2("ip.sh")
		return nt
	})

	m.Get("/all", func() string {
		nt := getEc2("ip_all.sh")
		return nt
	})

	m.Get("/urls", func() string {
		nt := getUrls("url")
		return nt
	})

	m.Get("/:name", func(params martini.Params) string {
		fileName := params["name"]
		fmt.Printf("fileName " + fileName)
		nt := getFileData(fileName)
		fmt.Printf("nt " + nt)
		return nt
	})

	m.Run()
}

func getEc2(fileName string) string {
	cmd := exec.Command("bash", fileName)
	cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}

func getUrls(fileName string) string {
	//cmd := exec.Command("cat", fileName)
	cmd := exec.Command("cat", "urls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}

func getFileData(fileName string) string {
	fmt.Println("getFileData inside " + fileName)
	cmd := exec.Command("cat", fileName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	} else {
		return string(out)
	}
	return "hello"
}
