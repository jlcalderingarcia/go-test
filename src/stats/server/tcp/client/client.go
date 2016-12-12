// Package to implement a simple TCP client to write text to the output
package client

import (
	"net"
	"strconv"
	"fmt"
)

//Structure to save the data of a simple tcp client
type OutTCPClient struct {
	Server string
	Port int
	Opened bool
	conn *net.TCPConn
}

//Function to create a new simple tcp client
func New(server string, port int) *OutTCPClient {
	res := new(OutTCPClient)
	res.Server = server
	res.Port = port
	res.Opened = false
	res.conn = nil
	return res
}

//Function to open the tcp client connection
func (c *OutTCPClient) Open() {
	//If the client is not oppened
	if !c.Opened && c.conn == nil {		
		//Find the server address
		tcpAddr, err := net.ResolveTCPAddr("tcp", c.Server + ":" + strconv.Itoa(c.Port))
		fmt.Printf("Successfully connected to: %s:%d\n", c.Server, c.Port)
		
		//If the address is valid
		if err == nil {
			//Try to open a connection
			c.conn, err = net.DialTCP("tcp", nil, tcpAddr)
			
			//If the connection succed, mark the client as opened
			if err == nil {
				c.Opened = true
			} else {
				fmt.Printf("Error: %s\n", err.Error())
			}
		} else {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}

//Function to close the tcp client connection
func (c *OutTCPClient) Close() {
	//Check if the client is connected
	if c.Opened && c.conn != nil {
		//Close the connection
		c.conn.Close()
	}
	
	//Clean the connection reference
	c.conn = nil
	//Set the client as closed
	c.Opened = false
}

//Function to write bytes out to the tcp connection
func (c *OutTCPClient) Write(bytes []byte){
	//If the client is connected and the connection opened
	if c.Opened && c.conn != nil {
		//Write the bytes to the server
		c.conn.Write(bytes)
	}
}
