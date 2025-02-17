/* TODO:Automatic server search..?*/
/* TODO:Better UI (at least erase printed lines, at best stop printIncomingMessages() while user is printing or something idk*/
/* TODO:Arguments for nick & server address*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func printIncomingMessages(connection net.Conn) {
	for {
		reader := bufio.NewReader(connection)
		message, _ := reader.ReadString('\n')
		fmt.Print(message)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("net.Dial ERROR: ", err)
	}
	fmt.Print("Write your nickname: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	fmt.Fprint(conn, text)

	go printIncomingMessages(conn)

	for {
		text, _ := reader.ReadString('\n')
		fmt.Fprint(conn, text)
	}
}
