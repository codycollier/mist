// partyline implements a basic chat server in the style of a party phone line
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// chatUser is a remotely connected party line chat participant.  Each user has
// an input loop running in a go routine, which takes messages typed by the
// user and sends them to the party line.
type chatUser struct {
	reader    *bufio.Reader
	writer    *bufio.Writer
	nick      string
	partyline *partyLine
}

func (user *chatUser) initialize(connection net.Conn, partyline *partyLine) {

	user.reader = bufio.NewReader(connection)
	user.writer = bufio.NewWriter(connection)
	user.partyline = partyline

	user.setNick()
	welcome := fmt.Sprintf("Welcome %s! Joining the line...\n", user.nick)
	user.send(welcome)
	user.partyline.addUser(user)

	go user.inputLoop()

}

func (user *chatUser) setNick() {
	user.writer.WriteString("Enter Nickname: ")
	user.writer.Flush()
	nick, _ := user.reader.ReadString('\n')
	user.nick = strings.TrimSpace(nick)
}

func (user *chatUser) inputLoop() {
	for {

		// block and wait for new input from the user
		line, _ := user.reader.ReadString('\n')

		// send input to the partyline
		msg := new(message)
		msg.user = user
		msg.text = line
		user.partyline.input <- msg
	}
}

func (user *chatUser) send(msg string) {
	user.writer.WriteString(msg)
	user.writer.Flush()
}

// message represents the details of a singe chat message from a user
type message struct {
	text string
	user *chatUser
}

// partyLine is a service which takes messages on an input channel and
// distributes those messages to all participants on the line.
type partyLine struct {
	users []*chatUser
	input chan *message
}

func (p *partyLine) start() {
	p.users = make([]*chatUser, 0)
	p.input = make(chan *message)
	go p.service()
}

func (p *partyLine) service() {
	for {
		select {
		case msg := <-p.input:
			log.Printf("%s> %s", msg.user.nick, msg.text)
		}
	}
}

func (p *partyLine) addUser(user *chatUser) {
	p.users = append(p.users, user)
}

// handleConnection initializes new incoming users
func newUser(connection net.Conn, partyline *partyLine) {
	log.Printf("New connection from: %s\n", connection.RemoteAddr())
	user := new(chatUser)
	user.initialize(connection, partyline)
	log.Printf("New user online %s@%s", user.nick, connection.RemoteAddr())
}

// Listen, accept connections, and connect them to the party line
func main() {

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	log.Println("Starting up...")

	partyline := new(partyLine)
	partyline.start()

	listener, _ := net.Listen("tcp", ":2323")
	defer listener.Close()

	for {
		connection, _ := listener.Accept()
		defer connection.Close()
		go newUser(connection, partyline)
	}

}
