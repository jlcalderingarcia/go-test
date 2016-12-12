// Package to define the TCP input server to save the word and chars stats
package tcp

import (
	"fmt"
	"testing"
	"stats/char"
	"stats/word"
	"os"
	"utils/random"
	"stats/server/tcp/client"
	"time"
)


//Function to test the creation of a new TCP server
func TestNew(t *testing.T){
	//Create the underlying data structures
	cstats := char.New()
	wstats := word.New()
	
	//Initialize a server on the given default port
	server, err := New(5555, cstats, wstats)
	
	if err != nil {
		t.Error("Server creation throws and error")
	}
	
	if server.Port != 5555 {
		t.Error("The port wasn't set succesfully")
	}
	
	if server.CharStats != cstats {
		t.Error("The char stats wasn't set succesfully")
	}
	
	if server.WordStats != wstats {
		t.Error("The word stats wasn't set succesfully")
	}
}

//Test the listening of the TCP server
func TestListenTCPServer(t *testing.T){
	//Create the underlying data structures
	cstats := char.New()
	wstats := word.New()
	
	//Initialize the server
	server, err := New(5550, cstats, wstats)
	
	if err != nil {
		t.Error("Server creation throws and error")
	}
	
	//Start listen in the server
	go server.Listen()
	
	//Crate a channel to inform this goroutine the TCP client end inserting data
	channel := make(chan int)
	
	//Send random data to the server in another go routine
	go sendValidatableClientData(t, server, channel)

	//Waits for the TCP client to complete
	<- channel
}

//Create a simple TCP client, send random data to the server and validate the server output
func sendValidatableClientData(t *testing.T, server *InTCPServer, channel chan int){	
	//Generate random test data
	stream, words, wcounts, chars, ccounts:= randomtext.GenerateRandomTestValues(1000)

	//Create a simple TCP client to send data to the listening server
	client := client.New("localhost", 5550)
	client.Open()
	client.Write([]byte(stream))//Write the output stream
	client.Close()

	//Leave a second to the server to process the input recently sended	
	time.Sleep(time.Second)
	
	//Validate the words counting
	if server.WordStats.Count != uint64(len(words)) {
		t.Errorf("The server misscounted the words, it has %d but there must be %d", server.WordStats.Count, len(words))
	}
	
	//Validate the words order
	mywords := server.WordStats.GetFirstN(100)
	l := 100
	if len(mywords) < l {
		l = len(mywords)
	}
	for i := 0; i < l; i++ {
		if mywords[i] != words[i] && wcounts[mywords[i]] != wcounts[words[i]] {
			t.Errorf("Misscounted word: %s", words[i])
		}
	}
	
	//Validate the chars counting
	if server.CharStats.Count != uint8(len(chars)) {
		t.Errorf("The server misscounted the chars, it has %d but there must be %d", server.CharStats.Count, len(chars))
	}
	
	//Validate the character order
	mychars := server.CharStats.GetFirstN(10)
	l = 10
	if len(mychars) < l {
		l = len(mychars)
	}
	for i := 0; i < l; i++ {
		if chars[i] != mychars[i] && ccounts[chars[i]] != ccounts[mychars[i]] {
			t.Errorf("Misscounted char: %s", string([]uint8{ chars[i] }))
		}
	}
	
	//Signal this go routine is done
	channel <- 1
}

func TestParalellProcessing(t *testing.T){
	//Create the underlying data structures
	cstats := char.New()
	wstats := word.New()
	
	//Initialize the server
	server, err := New(5551, cstats, wstats)
	
	if err != nil {
		t.Error("Server creation throws and error")
	}
	
	//Start listen in the server
	go server.Listen()
	
	//Create a comunication channel to syncronize the goroutines
	channel := make(chan int)
	
	//Call 100 simultaneous goroutines
	for i := 0; i < 100; i++ {
		go sendSimpleClientData(t, server, channel)		
	}

	//Wait for the goroutines to complete
	for i := 0; i < 100; i++{
		<- channel
	}
}

func sendSimpleClientData(t *testing.T, server *InTCPServer, channel chan int) {
	//Generate random test data
	stream, _, _, _, _:= randomtext.GenerateRandomTestValues(100)

	//Create a simple TCP client to send data to the listening server
	client := client.New("localhost", 5551)
	client.Open()
	client.Write([]byte(stream))//Write the output stream
	client.Close()
	
	//Leave a secound to the server to process the requests
	time.Sleep(time.Second)
	
	//Signal the goroutine is done
	channel <- 1
}

//Check the error variable and print the error string in the err output
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	}
}