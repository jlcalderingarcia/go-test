// wordstats_test.go
package word

import (
	"testing"
	"math/rand"
	"time"
	"utils/trie"
	"utils/random"
)

//Test the initialization of WordStat 
func TestInit (t *testing.T){
	ws := New()
	if ws.Count != 0 {
		t.Error("The Count must be 0")
	}
	if ws.Root == nil {
		t.Error("The trie root must be an initial one")
	}
	if len(ws.Index) != 0 {
		t.Error("The index map must be empty")
	}
	if len(ws.Words) != 1024 {
		t.Error("The words list must have at leats 1024 words allocated")
	}
	if ws.Index == nil || len(ws.Index) != 0 {
		t.Error("The words index must be allocated and empty")
	}
}

//Test the word insertion on the data structure
func TestInsert (t *testing.T) {
	ws := New()
	words, _ := generateRandomTestValues(ws, 1000) //Generate random words and process them randomly
	
	if int(ws.Count) != len(words) { //Check the number of words match
		trie.Print(ws.Root)
		t.Errorf("There are %d counted words, but are %d words", ws.Count, len(words))
	}
}

//Test the order and extraction of the first N words in the structure
func TestGetFirstN (t *testing.T) {
	ws := New()
	words, counts := generateRandomTestValues(ws, 1000) //Generate random words and process then randomly

	//Extract the first n results	
	n := 100
	nwords := ws.GetFirstN(uint64(n))
	
	if len(nwords) != n {
		t.Errorf("The number of words is wrong, it is %d and must be %d", len(nwords), n)
	}
	
	for i := 0; i < len(nwords); i++ { //Check the words are in correct order (they might be in diferent order but have the same word count)
		if nwords[i] != words[i] && counts[words[i]] != counts[nwords[i]] {
			t.Errorf("There word %s is missplaced, it has %d repetitions", nwords[i], counts[nwords[i]])	
		}
	}
}

//Generate random words and process them randomly
func generateRandomTestValues (ws *WordStats, nwords uint) (words []string, counts map[string]uint64){
	rand.Seed(time.Now().UnixNano()) //Initialize the random seed with the current timestamp
	counts = make(map[string]uint64) //Create the word count map
	words = make([]string, nwords) //Create the words array
	
	//Generate random words
	for i := uint(0); i < nwords; i++ {
		word := randomtext.Word() //Generate random word
		_, exists := counts[word] //Check if the word already exists
		if exists {
			i-- //If the word already exists loop again
		} else {
			counts[word] = 0 //If the word doesn't exists, save it
			words[i] = word
		}
	}
	
	//Duplicate random words
	for i := uint(0); i < nwords * 20; i++ {
		pos := rand.Intn(len(words)) //Get a random word and process it
		counts[words[pos]]++ //Update the word count
		ws.Insert(words[pos]) //Insert the word in the data structure
	}
	
	//Sort the words by counting
	for i := 0; i < len(words) - 1; i++ {
		for j := i + 1; j < len(words); j++ {
			if counts[words[i]] < counts[words[j]] {
				words[i], words[j] = words[j], words[i]
			}
		}
	}
	
	return
}