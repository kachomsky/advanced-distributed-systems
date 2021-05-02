/*
Name: Alejandro Garcia Carballo
User: u188873
*/

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

func trimInput(text string) string {
	result := ""
	if runtime.GOOS == "windows" {
		result = strings.TrimRight(text, "\r\n")
	} else {
		result = strings.TrimRight(text, "\n")
	}
	return result
}

func server(port string, waitgroup *sync.WaitGroup) {
	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":"+port)
	// run loop forever (or until ctrl-c or write stop)
	for {
		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Port could be in use")
		} else {
			// will listen for message to process ending in newline (\n)
			message, err := bufio.NewReader(conn).ReadString('\n')

			if err != nil {
				break
			}

			// output message received
			fmt.Print(string(message))
			messageTrim := trimInput(message)
			messageTrim = strings.Split(messageTrim, ":")[1]
			if messageTrim == "stop" {
				os.Exit(0)
			}
		}

	}
}

func client(configs []string, id string, waitgroup *sync.WaitGroup) {
	fmt.Print("Text to send: ")
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// send to socket
		for _, element := range configs[1:] {
			config := strings.Split(element, ":")
			conn, err := net.Dial("tcp", config[0]+":"+config[1])
			if err != nil {
				fmt.Println("The server is not ready")
			} else {
				fmt.Fprintf(conn, "ID"+id+" sent:"+text+"\n")
			}

		}

		if trimInput(text) == "stop" {
			os.Exit(0)
		}
	}
}

func readConfig(path string) string {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {
	var waitgroup sync.WaitGroup
	file := os.Args[1]
	lines := strings.Split(readConfig(file), "\n")
	id := strings.Split(lines[0], ":")[2]
	waitgroup.Add(1)
	go server(strings.Split(lines[0], ":")[1], &waitgroup)
	waitgroup.Add(1)
	go client(lines, id, &waitgroup)
	waitgroup.Wait()
}
