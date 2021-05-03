/*
Name: Alejandro Garcia Carballo
User: u188873
Lab 2: Echo algorithm
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
	"time"
)

type Node struct {
	ip      string
	port    string
	receive bool
	send    bool
}

var currentNode Node
var currentId string
var neighbors []Node
var initiator bool
var initiatorMessageSent bool
var parent Node

func configNodes(path string) {
	content, err := ioutil.ReadFile(path)
	initiator = false

	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	firstNodeSlice := strings.Split(lines[0], ":")
	if len(firstNodeSlice) >= 4 {
		if firstNodeSlice[3] == "*" {
			initiator = true
		}
	}
	ipCurrent := firstNodeSlice[0]
	portCurrent := firstNodeSlice[1]
	currentId = firstNodeSlice[2]

	currentNode = Node{ipCurrent, portCurrent, false, false}

	for _, element := range lines[1:] {
		neighbourConf := strings.Split(element, ":")
		neighbors = append(neighbors, Node{neighbourConf[0], neighbourConf[1], false, false})
	}
}

func server(s Node) {
	fmt.Println("Launching server...")
	ln, _ := net.Listen("tcp", s.ip+":"+s.port)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n')
		if string(message) != "" {
			msgSplit := strings.Split(message, ":")
			fmt.Println("Message received from node ID: ", msgSplit[len(msgSplit)-1])
			updateReceived(message, msgSplit[0], msgSplit[1])
		}

	}

}

func sendMessage(s string, n Node) {
	conn, _ := net.Dial("tcp", n.ip+":"+n.port)

	defer conn.Close()
	fmt.Fprintf(conn, currentNode.ip+":"+currentNode.port+":"+currentId+":"+s)
	fmt.Printf("Message sent to %s:%s \n", n.ip, n.port)
}

func checkNeighborServer(n []Node) bool {
	for i := 0; i < len(n); i++ {
		for {
			conn, err := net.Dial("tcp", n[i].ip+":"+n[i].port)
			fmt.Println("Looking for " + n[i].ip + ":" + n[i].port)
			time.Sleep(3000 * time.Millisecond)
			if err == nil {
				conn.Close()
				break
			}
		}
	}

	return true
}

func sendMessageAllNeighbors(message string) {
	for i := 0; i < len(neighbors); i++ {
		if neighbors[i].port != parent.port {
			sendMessage(message, neighbors[i])
			neighbors[i].send = true
		}
	}
}

func trimInput(text string) string {
	result := ""
	if runtime.GOOS == "windows" {
		result = strings.TrimRight(text, "\r\n")
	} else {
		result = strings.TrimRight(text, "\n")
	}
	return result
}

func findNodeByIdPort(ip string, port string) int {
	index := -1
	for i := 0; i < len(neighbors); i++ {
		if neighbors[i].ip == ip && neighbors[i].port == port {
			index = i
		}
	}
	return index
}

func updateReceived(message string, ip string, port string) {
	posNode := findNodeByIdPort(ip, port)
	n := neighbors[posNode]
	if initiator {
		neighbors[posNode].receive = true
	} else {
		if parent.ip == "" {
			parent = n
			neighbors[posNode].receive = true
			sendMessageAllNeighbors(currentId)
		} else {
			neighbors[posNode].receive = true
		}
	}
}

func allNeighboursReceived() bool {
	allReceived := true
	for _, node := range neighbors {
		if parent.port != node.port {
			if !node.receive {
				allReceived = false
			}
		}
	}
	return allReceived
}

func main() {
	file := os.Args[1]
	configNodes(file)
	fmt.Println("This node is ", currentNode.ip, currentNode.port)
	fmt.Println("ID: ", currentId)
	if initiator {
		fmt.Println("This node is an initiator")
	}
	fmt.Println("\n ====   Neighbours for this node   ====\n")
	for i := 0; i < len(neighbors); i++ {
		fmt.Println(neighbors[i].ip + ":" + neighbors[i].port)
	}
	fmt.Println("\n======================================\n")

	parent = Node{"", "", false, false}
	go server(currentNode)
	if checkNeighborServer(neighbors) {
		finish := false
		for {
			if initiator && !initiatorMessageSent {
				fmt.Println("Initiator " + currentId + ": sending message")
				sendMessageAllNeighbors(currentId)
				initiatorMessageSent = true
			} else {
				time.Sleep(3000 * time.Millisecond)
				if allNeighboursReceived() {
					if initiator {
						fmt.Println("\nDone ")
						finish = true
					} else {
						sendMessage(currentId, parent)
						fmt.Println("\nDone ")
						finish = true
					}
				}
			}
			if finish {
				break
			}
		}
	}
}
