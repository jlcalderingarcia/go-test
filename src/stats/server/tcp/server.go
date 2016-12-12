// Package to define the TCP input server to save the word and chars stats
package tcp

import (
	"stats/char"
	"stats/word"
	"strconv"
	"net"
	"os"
	"strings"
	"io"
	"fmt"
)

type InTCPServer struct {
	CharStats *char.CharStats
	WordStats *word.WordStats
	Port int
	wordsChannel chan string
	charsChannel chan uint8
}

//Create a new input TCP Server
func New(port int, charStats *char.CharStats, wordStats *word.WordStats) (server *InTCPServer, err error){
	server = new(InTCPServer)
	server.Port = port
	server.CharStats = charStats
	server.WordStats = wordStats
	server.wordsChannel = make(chan string)
	server.charsChannel = make(chan uint8)
	return
}

func (s *InTCPServer) Listen() {
	fmt.Printf("TCP interface listening in: :%d\n", s.Port)

	//Initialize the listening port
	tcp, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		fmt.Println("Error initializing the listener: ", err.Error())
		os.Exit(1)
	}
	// Close the listener on the function ends
	defer tcp.Close()
	
	//Start listening for words
	go s.proccessWord()
	
	//Start listening for chars
	go s.processChar()
	
	for {
		// Wait for the clients to connect
		conn, err := tcp.Accept()
		if err != nil {
			continue
		}

		// Handle the incoming connection.
		go s.handleTCPRequests(conn)
	}
}

//This function handle the tcp conection
func (s *InTCPServer)handleTCPRequests(conn net.Conn) {	
	//Close the connection when the function ends
	defer conn.Close()
	
	//Create the buffer to listen from the client
	buffer := make([]byte, 1024 * 4)
	
	//Create the string to hold the remaining text
	line := ""
	
	//Set the valid characters limit
	a, z := "a"[0], "z"[0]

	//Position of the initial word index and the current index
	i, j := 0, 0 
			
	for {		
		//Read data from the buffer
		n, err := conn.Read(buffer)
		
		if err  != nil {
			//If the stream ends and there's a word to process, process it
			if err == io.EOF && len(line) > 0 {
				s.wordsChannel <- line
			}
			break //Breaks the loop to close the conection
		} else {
			//Convert the readed bytes to lowercase string and joins it to the remaining text from the last read
			line = line + strings.ToLower(string(buffer[0:n]))
			//Start processing the chars and words on the readed input
			for i = 0; i < len(line) && j < len(line); j++ {
				//Validate the current character is valid
				if line[j] < a || line[j] > z {
					//If there's a word readed send it to the channel to be processed
					if i < j {
						s.wordsChannel <- line[i:j]
					}
					//Clear the word initial index and set it to the next index
					i = j + 1
				} else {
					//Process send the current character to the channel to be processed
					s.charsChannel <- line[j]
				}
			}
			//Validate the new extart
			if i > len(line) {
				i = len(line)
			}
			//Keep only the unprocesed text
			line = line[i:]
			
			//Set the starting scanning point
			j = len(line)
		}
	}
}

//This function handle the words processing
func (s *InTCPServer)proccessWord(){
	for {
		//Wait for a word in the channel
		word := <-s.wordsChannel
		
		//Process the word
		s.WordStats.Insert(word)
	}
}

//This function handle the chars processing
func (s *InTCPServer) processChar(){
	for {
		//Wait for a char in the channel
		char := <-s.charsChannel
		
		//Process the char
		s.CharStats.Listen(char)
	}
}