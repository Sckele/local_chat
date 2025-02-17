/* TODO:telegram connection?*/
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func readMessageFromUserConnection(user user) (string, error) {
	message := make([]byte, 1024)
	message_length, err := user.address.Read(message)
	if err != nil {
		return "", err
	}

	formatted_message := strings.Trim(string(message[:message_length]), "\n")
	fmt.Println("Received message: ", "\""+formatted_message+"\"", "from: ", user.nick)
	return formatted_message, nil
}
func sendMessageToUser(user user, message string) (err error) {
	if message == "" {
		fmt.Println("Message is empty!")
		return nil
	}
	_, err = fmt.Fprintln(user.address, message)
	if err != nil {
		return err
	}
	fmt.Println("Sent message:", "\""+message+"\"", "to: ", user.nick)
	return nil
}

func broadcastMessages(users *[]*user, messages <-chan string) {
	var message string
	for {
		message = <-messages
		if message == "" {
			continue
		}
		for _, user := range *users {
			sendMessageToUser(*user, message)
		}
	}
}

func handleUserConnection(user *user, messages chan<- string) {
	defer (*user).address.Close()
	connectionData := (*user).address.LocalAddr().String() + " " + (*user).address.LocalAddr().Network() + " " + (*user).address.RemoteAddr().String()
	fmt.Println("New connection: ", connectionData)

	message, err := readMessageFromUserConnection(*user)
	if err != nil {
		fmt.Println(err, connectionData)
		return
	}
	(*user).nick = message

	for {
		message, err = readMessageFromUserConnection(*user)
		if err != nil {
			fmt.Println(err, connectionData)
			return
		}
		message = strings.Join([]string{"[", (*user).nick, "] ", message}, "")
		messages <- message
	}
}

func main() {
	const port string = "8000"
	address, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("Listening on port " + port)

	messages := make(chan string, 10)
	var users []*user
	go broadcastMessages(&users, messages)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		user := user{address: connection, nick: "-"}
		users = append(users, &user)
		go handleUserConnection(&user, messages)
	}
}
