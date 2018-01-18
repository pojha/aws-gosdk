package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, getEc2("ip.sh"))
}
func handlerAll(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, getEc2("ip_all.sh"))
}
func handlerFile(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, getFileData("url"))
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/all", handlerAll)
	http.HandleFunc("/url", handlerFile)
	http.ListenAndServe(":3000", nil)
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
	//  fmt.Printf("Output \n",out.String());
	return out.String()
}

func getFileData(fileName string) string {
	//cmd := exec.Command("cat", fileName)
	cmd := exec.Command("cat", "urls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return string(out)
}
