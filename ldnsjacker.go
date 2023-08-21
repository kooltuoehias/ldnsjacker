package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
)

var (
	pathEtcHosts string = "C:/Windows/System32/drivers/etc/hosts"
	localhost    string = "127.0.0.1"
	markerOfJack string = "###LJACK START###"
	quit         string = "q"
)

func main() {

	if runtime.GOOS != "windows" {
		pathEtcHosts = "/etc/hosts"
	}

	restore(pathEtcHosts)
	fmt.Println("OS: ", runtime.GOOS)
	fmt.Println("Path: ", pathEtcHosts)
	defer restore(pathEtcHosts)
	add(pathEtcHosts, os.Args[1:])

	pressQuit()
}

func add(hostsFilePath string, domains []string) {
	// Open the hosts file for writing
	f, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening hosts file:", err)
		return
	}
	defer f.Close()

	// Append the new records to the hosts file
	_, err = f.WriteString(markerOfJack + "\n")
	if err != nil {
		fmt.Println("Error writing to hosts file:", err)
		return
	}
	for _, domain := range domains {
		_, err = f.WriteString(fmt.Sprintf("\n%s %s\n", localhost, domain))
		if err != nil {
			fmt.Println("Error writing to hosts file:", err)
			return
		}

	}

}

func restore(hostsFilePath string) {

	f, err := os.Open(hostsFilePath)
	if err != nil {
		fmt.Println("Error opening hosts file:", err)
	}
	defer f.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(f)
	writeback := true
	for scanner.Scan() {
		text := scanner.Text()
		if text == markerOfJack {
			writeback = false
		}
		if writeback {
			_, err = buf.WriteString(text + "\n")
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	err = os.WriteFile(hostsFilePath, buf.Bytes(), 0666)

	if err != nil {
		fmt.Println("Error writing to hosts file:", err)
	}
}

func pressQuit() {
	fmt.Printf("Type %s to restore and quit: \n", quit)
	for reader := bufio.NewReader(os.Stdin); ; {
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(text) == quit {
			return
		}
	}

}
