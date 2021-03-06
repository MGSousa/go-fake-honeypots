package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	PORT     = 23 // TODO: make this as a flag OR build with -X to change var
	hostname = "GW-EXTERNAL"
)

func floatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	basePath := filepath.Dir(ex)

	f, err := os.OpenFile(fmt.Sprintf("%s/fh-telnet.log", basePath), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	server, _ := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if server == nil {
		panic("couldn't start listening: ")
	}
	conns := clientConns(server)
	for {
		go handleConn(<-conns)
	}
}

func clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, _ := listener.Accept()
			if client == nil {
				fmt.Printf("couldn't accept: ")
				log.Println("couldn't accept: ")
				continue
			}
			i++

			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			log.Printf("CLIENT OPENED CONNECTION %d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			ch <- client
		}
	}()
	return ch
}

func ping(command string, client net.Conn) {
	cmd := command
	if strings.Contains(cmd, "ping") {
		result := strings.Split(cmd, " ")
		host := result[1]

		fmt.Println(result[0])
		fmt.Println(result[1])
		fmt.Printf("PING %v (%[1]v) 56(84) bytes of data.\n", host)

		for i := 0; i < 4; i++ {
			time.Sleep(2 * time.Second)
			time_ping := floatToString((rand.Float64() * 10) + 20)
			ttl := strconv.Itoa((rand.Intn(5) * 10))
			fmt.Printf("64 bytes from %v: icmp_seq=%v ttl=%v time=%.2v ms\n", host, i, ttl, time_ping)
			stringArray := []string{"64 bytes from ", host, ": icmp_seq =", strconv.Itoa(i), " ttl =", ttl, " time=", time_ping[0:3], " ms \n"}
			s := strings.Join(stringArray, " ")
			client.Write([]byte(s))

		}
	}

}

func handleConn(client net.Conn) {
	ena := "nope"
	conft := "nope"
	b := bufio.NewReader(client)
	c := bufio.NewReader(client)
	motd := "This device is for authorized personnel only. \n" +
		"If you have not been provided with permission to \n" +
		"access this device - disconnect at once. \n" +
		"*** Login Required.  Unauthorized use is prohibited *** \n" +
		"*** Ensure that you update the system configuration *** \n" +
		"*** documentation after making system changes.      *** \n" +
		"User Access Verification: "

	log.Println("Sending MOTD to client")
	client.Write([]byte(motd))
	user, _ := b.ReadBytes('\n')
	client.Write([]byte("Password: "))
	pass, _ := b.ReadBytes('\n')

	fmt.Println("user : ", string(user))
	fmt.Println("pass :", string(pass))
	log.Printf("user : %v ", string(user))
	log.Printf("pass : %v", string(pass))

	for {
		line, err := c.ReadBytes('\n')
		if err != nil {
			break
		}
		stringa := string(line)
		cmd := strings.Trim(stringa, " \r\n")
		fmt.Println(cmd)

		if cmd == "ena" || cmd == "enab" || cmd == "enabl" || cmd == "enable" || cmd == "sudo su" {
			ena = "yes"
			fmt.Println("ENA ", ena)
		}
		if cmd == "configuration terminal" || cmd == "configure terminal" || cmd == "conf termi" || cmd == "conf t" {
			conft = "yes"
			fmt.Println("Confetti", conft)
		}

		if ena == "yes" && conft == "nope" {
			fmt.Println("ENABLE USER")
			ping(cmd, client)
			stringArray1 := []string{hostname, "#"}
			s1 := strings.Join(stringArray1, " ")
			client.Write([]byte(s1))
		}

		if conft == "nope" && ena == "nope" {
			fmt.Println("ENABLE USER")
			stringArray3 := []string{hostname, ">"}
			s3 := strings.Join(stringArray3, " ")
			client.Write([]byte(s3))
			ping(cmd, client)
		}

		if conft == "yes" && ena == "nope" {
			fmt.Println("ENABLE USER")
			client.Write([]byte("You do not have sufficient privileges to execute this command"))
			conft = "nope"
		}

		if conft == "yes" && ena == "yes" {
			fmt.Println("ENABLE USER")
			ping(cmd, client)
			stringArray := []string{hostname, "(config)#"}
			s := strings.Join(stringArray, " ")
			client.Write([]byte(s))
		}

		fmt.Println("CMD: ", cmd)
		log.Println("CMD Sent by client:", cmd)

		if cmd == "exit" && ena == "nope" {
			fmt.Printf("CLOSING conn")
			client.Write([]byte("Bye..\n"))
			client.Close()
		}
		if cmd == "exit" && ena == "yes" {
			ena = "nope"
			conft = "nope"
			cmd = ""
		}

	}
}
