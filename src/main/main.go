// Word stats cmd utility
package main

import (
	"fmt"
	"os"
	"strconv"
	"stats/server/http"
	"stats/server/tcp"
	"stats/char"
	"stats/word"
)

func main() {
	//Initialize the variables
	var (
		wordsPortNumber int64
		statsPortNumber int64
		err error
	)
	
	//Initialize the default listen port
	wordsPortNumber = 5555 
	statsPortNumber = 8080
	
	//Analyze the program arguments to find the flags to set up the ports
	for i := 0; i < len(os.Args) - 1; i++{ 
		switch os.Args[i]{
			case "-w":
				wordsPortNumber, err = strconv.ParseInt(os.Args[i + 1], 10, 32)
				if err != nil || wordsPortNumber <= 0 {
					printUsage()
				}
				i++
			case "-s":
				statsPortNumber, err = strconv.ParseInt(os.Args[i + 1], 10, 32)
				if err != nil || statsPortNumber <= 0{
					printUsage()
				}
				i++
		}
	}
	fmt.Printf("Words port: %d\n", wordsPortNumber)
	fmt.Printf("Stats port: %d\n", statsPortNumber)
	
	//Initialize the base data structures
	cstats := char.New()
	wstats := word.New()
	
	//Create the HTTP stats server and start listening in other goroutine
	httpServer, _ := http.New("/stats", int(statsPortNumber), cstats, wstats)
	go httpServer.Listen()
	
	//Create the TCP listener sever and start listening on this goroutine
	tcpServer, _ := tcp.New(int(wordsPortNumber), cstats, wstats)
	tcpServer.Listen()
}

//Print the progam usage message
func printUsage(){
	fmt.Println("Usage: %s -w <word server port> -s <stats server port>")
}

//Validate the given error variable and output the error message in the standar error stream
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
