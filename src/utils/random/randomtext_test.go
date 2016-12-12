// Package used to generate random chars and words
package randomtext

import (
	"time"
	"math/rand"
	"testing"
	"utils/random"
)

//Test the char random generation
func TestChar(t *testing.T){
	//Initialize the random seed with the current timestamp
	rand.Seed(time.Now().UnixNano())
	
	//Generate 10000 characters randomly and count them
	m := make(map[uint8]uint64, 0)
	for i := 0; i < 10000; i++ {
		c := randomtext.Char()
		m[c]++
	}
	
	//All chracates must be in the map (the probabilities they don't are minimal)
	chars := "az"
	if int(chars[1] - chars[0]) > len(m) {
		t.Error("There are missing characters")
	}
	
	//Verify all the characters are in the appropiate range
	for c, count := range m {
		if count == 0 {
			t.Error("There are characters no generated")
		}
		if c < chars[0] || c > chars[1] {
			t.Error("There are invalid characters")
		}
	}
}


//Test the random word generation
func TestWord(t *testing.T){
	//Initialize the random seed with the current timestamp
	rand.Seed(time.Now().UnixNano())
	
	//Generate 1000 words randomly, if some repeats more than twice, they ae not random
	count := 0
	for i := 0; i < 1000; i++ {
		//Generate two random words and compare thems
		if randomtext.Word() == randomtext.Word() {
			count++
			if count > 2 {
				t.Error("The generated words are not random")
			}
		}
	}
}

//Test the random test value generation
func TestGenerateRandomTestValues(t *testing.T){
	//Initialize the random seed with the current timestamp
	rand.Seed(time.Now().UnixNano())
	
	//Generate a random number betwen 100 and 999
	nwords := rand.Intn(900) + 100
	
	//Generate the test input
	stream, words, wcounts, chars, ccounts := GenerateRandomTestValues(uint(nwords))
	
	if len(stream) == 0 || len(words) == 0 || len(wcounts) == 0 || len(chars) == 0 || len(ccounts) == 0 {
		t.Error("The function doesn't generate any output stream")
	}
	
	if len(words) != len(wcounts) {
		t.Error("The function misscounted the words")
	}
	
	if len(chars) != len(ccounts) {
		t.Error("The function misscounted the chars")
	}
	
	if len(words) != nwords {
		t.Error("The function didn't generate the correct words count")
	}
}