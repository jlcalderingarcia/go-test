// Package to manage the char stats generateds by the app
package char

import (
	"testing"
	"utils/random"
	"math/rand"
	"time"
)

//Test the new initialized chars stats structure has the proper attributes
func TestInit (t *testing.T){
	count := int("z"[0] - "a"[0] + 1)
	
	cs := New() //Initialize the data structure
	if cs.Count != 0 {
		t.Error("The initial characters count must be 0")
	}
	if len(cs.Index) != 0 {
		t.Error("The initial map with the characters sort vector index must be empty")
	}
	if len(cs.Counters) != 0 {
		t.Error("The initial counters map must be empty")
	}
	if len(cs.Chars) != count {
		t.Error("The initial sorted vector must be big enough to contain all the characters")
	}
}

//Test the listen method of the data structure
func TestListen(t *testing.T){
	cs := New() //Initialize the data structure
	
	stream, counts := generateRandomCharacters(10000) //Genrate random input for the data structure
	
	//Make the data structure listen the randomly generated input
	for i := 0; i < len(stream); i++ {
		cs.Listen(stream[i])
	}
	
	if cs.Count == 0 {
		t.Error("There must be several characters in the data structure")
	}
	
	if int(cs.Count) != len(counts) {
		t.Errorf("There are %d distinct characters in the data structure when here must be %d", cs.Count, len(counts))
	}
	
	//Validate the counter of every character
	for char, count := range counts {
		if cs.Counters[char] != count {
			t.Errorf("The char %s has being misscounted, it has %d occurrences when it must have %d", string([]uint8{ char }), cs.Counters[char], count)
		}
	}
}

//Test the sorted characters structure
func TestGetFirstN(t *testing.T){
	cs := New() //Initialize de data structure
	
	stream, counts := generateRandomCharacters(10000) //Generate random input for the data structure
	
	//Make the data structure listen the randomly generated input
	for i := 0; i < len(stream); i++ {
		cs.Listen(stream[i])
	}
	
	//Extract the character random input
	chars := make([]uint8, len(counts))
	i := 0
	for c, _ := range counts {
		chars[i] = c
		i++
	}
	
	//Sort the chars in the random input by count
	for i = 0; i < len(chars) - 1; i++ {
		for j := i + 1; j < len(chars); j++ {
			if counts[chars[i]] < counts[chars[j]] {
				chars[i], chars[j] = chars[j], chars[i]
			}
		}
	}
	
	sorted := cs.GetFirstN(20) //Extract the first 20 characters most listeneds in the data structure
	//Validate the characters match or have the same counter
	for i = 0; i < 20; i++ {
		if sorted[i] != chars[i] && counts[sorted[i]] != counts[chars[i]] {
			t.Errorf("The character %s is missplaced, it has %d repetitions, but the character %s has %d and must be before", string([]uint8{sorted[i]}), counts[sorted[i]], string([]uint8{ chars[i] }), counts[chars[i]])
		}
	}
}

//Generate random input for the tests
func generateRandomCharacters(n uint64) (buffer []uint8, count map[uint8]uint64){
	rand.Seed(time.Now().UnixNano()) //Initialize the random seed with the current timestamp
	
	buffer = make([]uint8, n) //Buffer to save the generated character stream
	count = make(map[uint8]uint64, 0) //Map to count the repetitions of the caracters
	
	//Generate and count random characters
	for i := uint64(0); i < n; i++ {
		buffer[i] = randomtext.Char() //Generate a random char
		_, exists := count[buffer[i]] //Verified if the char has olready counted
		if !exists { //Insert it in the counters map if not
			count[buffer[i]] = 1
		} else {
			count[buffer[i]]++ //Count the new occurence
		}
	}
	
	return
}