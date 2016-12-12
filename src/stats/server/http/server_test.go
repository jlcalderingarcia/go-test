// out_test
package http

import (
	"fmt"
	"testing"
	"stats/char"
	"stats/word"
	"net/http"
	"io/ioutil"
)


func TestNew(t *testing.T){
	cstats := char.New()
	wstats := word.New()
	
	server, err := New("/stats", 8080, cstats, wstats)
	
	if err != nil {
		t.Error("Server creation throws and error")
	}
	
	if server.Route != "/stats" {
		t.Error("The route wasn't set succesfully")
	}
	
	if server.Port != 8080 {
		t.Error("The port wasn't set succesfully")
	}
	
	if server.CharStats != cstats {
		t.Error("The char stats wasn't set succesfully")
	}
	
	if server.WordStats != wstats {
		t.Error("The word stats wasn't set succesfully")
	}
}

func TestListen(t *testing.T){
	cstats := char.New()
	wstats := word.New()
	
	server, _ := New("/stats", 8080, cstats, wstats)
	
	go server.Listen()
	
	resp, err := http.Get("http://localhost:8080/stats?N=3")
	if err != nil {
		t.Error("Server isn't listening on the given port")
	} else {
		//defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(body)
		} else {
			t.Error("Cant read the server response")
		}
	}
}