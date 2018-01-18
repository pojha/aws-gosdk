// Writing files in Go follows similar patterns to the
// ones we saw earlier for reading.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"io"
	"os/exec"

)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	getEc2("epen_ip1")



	// For more granular writes, open a file for writing.
	f, err := os.Create("./dat2")
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	// You can `Write` byte slices as you'd expect.
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	// A `WriteString` is also available.
	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()

	// `bufio` provides buffered writers in addition
	// to the buffered readers we saw earlier.
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)

	// Use `Flush` to ensure all buffered operations have
	// been applied to the underlying writer.
	w.Flush()

}

func getEc2(fileName string) string {
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
	outStr := string(stdoutBuf.Bytes())

	err = ioutil.WriteFile(fileName, stdoutBuf.Bytes(), 0644)
	return outStr
}