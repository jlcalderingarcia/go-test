// Package to implement a simple TCP client to write text to the output
package client

import (
	"testing"
	"strconv"
	"time"
	"net"
	"io/ioutil"
)

//Validate the TCP client initialize correctly
func TestNew(t *testing.T){
	client := New("localhost", 5555)
	
	if client.Server != "localhost" {
		t.Error("The server hasn't being initialized correctly")
	}
	
	if client.Port != 5555 {
		t.Error("The port wasn't initialized propertly")
	}
	
	if client.Opened {
		t.Error("The opened indicator wasn't initialized propertly")
	}
	
	if client.conn != nil {
		t.Error("The TCP connection wasn't initialized propertly")
	}
}

func TestOpenClose(t *testing.T){
	//Create a simple TCP listener wich listen for just one connection and closes
	go createSimpleConnectionTCPListener(5555, make(chan string))
	
	//Wait 500 milliseconds for the server to open
	time.Sleep(time.Millisecond * 100)
	
	//Create the new client
	client := New("localhost", 5555)
	
	//Open the conection
	client.Open()
	
	//Wait for a 100 milliseconds
	time.Sleep(time.Millisecond * 100)
	
	if !client.Opened || client.conn == nil {
		t.Error("The client must be open now")
	}
	
	//Close the clinet
	client.Close()
	
	if client.Opened || client.conn != nil {
		t.Error("The client must be closed now")
	}
}

func TestWrite(t *testing.T){
	//Save the cannel
	channel := make(chan string)
	//Create a simple TCP listener wich listen for just one connection and closes
	go createSimpleConnectionTCPListener(5554, channel)
	
	//Wait 500 milliseconds for the server to open
	time.Sleep(time.Millisecond * 100)
	
	//Create the new client
	client := New("localhost", 5554)
	
	//Open the conection
	client.Open()
	
	//Wait for a 100 milliseconds
	time.Sleep(time.Millisecond * 100)
	
	//Write a text on the client
	text := "This are we testing the simple TCP client"
	client.Write([]byte(text))
	
	//Close the client
	client.Close()
	
	//Extract the text from the channel
	prove := <- channel
	
	if prove != text {
		t.Error("The TCP client isn't comunicating well")
	}
}

//Create a simple and only one connection listen TCP server that sends the received text over the channel provided
func createSimpleConnectionTCPListener(port int, output chan string){
	//Start listen on the TCP providen port
	ln, _ := net.Listen("tcp", ":" + strconv.Itoa(port))
	//Close the connection on end
	defer ln.Close()
	
	//Wait for an incoming connection
	conn, err := ln.Accept()
	if err == nil {
		//Close the connection on end
		defer conn.Close()
		
		//Read all the input
		input, _ := ioutil.ReadAll(conn)
		
		//Send the received content
		output <- string(input)
	}
}
