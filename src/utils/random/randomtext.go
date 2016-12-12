// Package used to generate random chars and words
package randomtext

import (
	"math/rand"
	"time"
)

//Generate a random word with a radom length betwen 1 and 100 and lower case text characters
func Word() string{
	word := make([]uint8, rand.Intn(100) + 1) //Generate the character base array
	for i := 0; i < len(word); i++ {
		word[i] = Char()//Generate every character
	}
	return string(word)
}

//Generate random lower case text characters [a..z]
func Char() uint8 {
	t := "az"
	return uint8(rand.Intn(int(t[1] - t[0])) + int(t[0]))
}

//Generate random words and process them randomly
func GenerateRandomTestValues (nwords uint) (stream string, words []string, wcounts map[string]uint64, chars []uint8, ccounts map[uint8]uint64){
	rand.Seed(time.Now().UnixNano()) //Initialize the random seed with the current timestamp
	words = make([]string, nwords) //Create the words array
	wcounts = make(map[string]uint64) //Create the word count map
	ccounts = make(map[uint8]uint64) //Create the chars count
	
	//Generate random words
	for i := uint(0); i < nwords; i++ {
		word := Word() //Generate random word
		_, exists := wcounts[word] //Check if the word already exists
		if exists {
			i-- //If the word already exists loop again
		} else {
			wcounts[word] = 0 //If the word doesn't exists, save it
			words[i] = word
		}
	}
	
	stream = ""
	//Duplicate random words
	for i := uint(0); i < nwords * 20; i++ {
		pos := rand.Intn(len(words)) //Get a random word and process it
		wcounts[words[pos]]++ //Update the word count
		if i < nwords * 20 - 1{
			stream = stream + words[pos] + " ";
		} else {
			stream = stream + words[pos];
		}
		//Count the used characters
		for j := 0; j < len(words[pos]); j++ {
			_, exists := ccounts[words[pos][j]]
			if !exists {
				ccounts[words[pos][j]] = 1
			} else {
				ccounts[words[pos][j]]++
			}
		}
	}
	
	//Sort the words by counting
	for i := 0; i < len(words) - 1; i++ {
		for j := i + 1; j < len(words); j++ {
			if wcounts[words[i]] < wcounts[words[j]] {
				words[i], words[j] = words[j], words[i]
			}
		}
	}
	
	//Save the characters
	chars = make([]uint8, len(ccounts))
	i := 0
	for c, _ := range ccounts {
		chars[i] = c
		i++
	}
	
	//Sort the characters by counting
	for i := 0; i < len(chars) - 1; i++ {
		for j := i + 1; j < len(chars); j++ {
			if ccounts[chars[i]] < ccounts[chars[j]] {
				chars[i], chars[j] = chars[j], chars[i]
			}
		}
	}
	
	return
}