// Package to manage the word stats managed by the app
package word

import (
	"utils/trie"
	"sync"
)

//Create the data structure to keep all the information
type WordStats struct {
	Root *trie.Trie //The trie root
	Words []*trie.Trie //The sorted vector of words (represented by trie left)
	Index map[*trie.Trie]uint64 //Index of the words in the ordered vector
	Count uint64 //The numbers of words currently in the structure
	mutex sync.Mutex
}

//Initialize an empty WordStats structure
func New() *WordStats {
	res := new(WordStats)
	res.Root = trie.New()
	res.Words = make([]*trie.Trie, 1024)
	res.Index = make(map[*trie.Trie]uint64)
	res.Count = 0
	return res
}

//Insert a new word into the structure
func (s *WordStats) Insert(word string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	t := trie.Insert(s.Root, word) //Insert the word into the trie and get the trie left representing the word
	_, exists := s.Index[t] 
	//If it's an unknown word, insert it
	if !exists {
		//If the ordered vector is full, duplicate it size
		if s.Count == uint64(len(s.Words)) {
			arr := make([]*trie.Trie, len(s.Words) * 2)
			copy(arr[0:len(s.Words)], s.Words)
			s.Words = arr
		}
		s.Words[s.Count] = t //Save the word into the ordered vector
		s.Index[t] = s.Count //Update the index
		s.Count++ //Update the counter
	} else { //If it a known word, update the index
		i := s.Index[t] //Find the index of the word
		for ; i > 0 && s.Words[i - 1].Count < t.Count; i-- { //Move to the start the most repeated word
			s.Index[s.Words[i - 1]]++ //Down the previous word
			s.Index[t]-- //Up the current word
			s.Words[i - 1], s.Words[i] = s.Words[i], s.Words[i - 1] //Update the vector order
		}
	}
}


//Extract the first N most repeated words
func (s *WordStats) GetFirstN(n uint64) []string{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	count := n
	if count > s.Count { //Get the lower value, the given N or the words count
		count = s.Count
	}
	res := make([]string, count)//Create the result vector
	for i := uint64(0); i < count; i++ { //Extract all the words
		res[i] = trie.GetWord(s.Words[i]) //Convert the trie to string
	}
	return res
}
