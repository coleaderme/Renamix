package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

func decode(dir string) {
	// open MasterTable file to read
	rfile, err := os.Open("mtable.sys")
	if err != nil {
		panic(err)
	}
	defer rfile.Close()
	// store MasterTable content to data variable (buffer)
	data, err := io.ReadAll(rfile)
	if err != nil {
		panic(err)
	}
	// loop over lines, split by ":" to name:id variables
	for _, d := range strings.Split(string(data), "\n") {
		// skip line with 0 length / line with # commented
		if len(d) != 0 && string(d[0]) != "#" {
			split := strings.Split(d, ":")
			name := split[0]
			id := split[1]
			// real magic happens here
			// if hashed path not exists it will error, else successfully decoded!
			if err := os.Rename(dir+id, dir+name); err != nil {
				fmt.Println("[-] Decode: Skipping..", err)
			} else {
				fmt.Println("[+] Decoded:\n" + dir + id + "\n=>\n" + dir + name)
			}
		}
	}
}

func encode(dir string) {
	// open MasterTable file as wfile variable to append content
	wfile, err := os.OpenFile("mtable.sys", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer wfile.Close()
	// loop Files array
	for _, name := range listFiles(dir) {
		id := getuuid()
		// write name:id
		if _, err := wfile.WriteString(name + ":" + id + "\n"); err != nil {
			fmt.Println("[-] write", err)
		}
		// rename to new hashed name
		if err := os.Rename(dir+name, dir+id); err != nil {
			fmt.Println("[-] Encode: ", err)
		} else {
			fmt.Println("[+] Encoded:\n" + dir + name + "\n=>\n" + dir + id)
		}
	}
	// end of loop, append one decoration line
	if _, err := wfile.WriteString("#========================\n"); err != nil {
		fmt.Println("[-] write", err)
	}
}

func getuuid() string {
	newUUID := uuid.New() // UUID of len 16
	return hex.EncodeToString(newUUID[:])
}

func listFiles(dir string) []string {
	// Read the directory
	items, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	files := make([]string, 0, 15) // static memory for 15 files ? why
	// loop items
	for _, i := range items {
		if !i.IsDir() {
			name := i.Name()
			fmt.Println("[+] " + name)
			files = append(files, name)
		}
	}
	return files
}

func help() {
	fmt.Println("USAGE: " + os.Args[0] + " /path/to/dir/ " + "[e|d]")
}

func main() {
	// dir := "example/path/to"
	if len(os.Args) < 3 {
		help()
		return
	}
	dir := os.Args[1]
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	operation := os.Args[2]

	if operation == "e" || operation == "E" {
		encode(dir)
	} else if operation == "d" || operation == "D" {
		decode(dir)
	} else {
		help()
	}
}
