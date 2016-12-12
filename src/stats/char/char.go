// Package to manage the char stats generateds by the app
package char

import (
	"sync"
)

//Structure to keep the data of the repetitions of the chars
type CharStats struct {
	Chars []uint8 //The ordered chars
	Index map[uint8]int //The map with the sorted index of the given char
	Counters map[uint8]uint64 //The map with the count of every char repetition
	Count uint8 //The count of chars used
	Listening chan int
	mutex sync.Mutex
}

//Initialize an empty data structure to manage the chars stats of the app
func New() *CharStats {
	count := int("z"[0] - "a"[0] + 1)
	res := new(CharStats)
	res.Count = 0
	res.Index = make(map[uint8]int, count)
	res.Counters = make(map[uint8]uint64, count)
	res.Chars = make([]uint8, count)
	res.Listening = make(chan int, 0)
	return res
}

//Listen a new chart to keep the stats
func (x *CharStats) Listen(char uint8) {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	
	index, exists := x.Index[char] //Checks if the char is already inserted
	if !exists { //If the character hasn't being inserted
		x.Chars[x.Count] = char //Save the chart in the ordered vector
		x.Counters[char] = 1 //Save the char repetition count
		x.Index[char] = int(x.Count) //Save the ordered index of the char
		x.Count++ //Update the characters count
	} else {
		x.Counters[char]++ //Update the given character count
		
		for ; index > 0 && x.Counters[x.Chars[index - 1]] < x.Counters[char]; index-- { //Bubble up the character
			x.Index[x.Chars[index - 1]]++ //Update the previos character index
			x.Index[char]-- //Update the given character inde
			x.Chars[index - 1], x.Chars[index] = x.Chars[index], x.Chars[index - 1] //Swap the previous and current character
		}
	}
}

//Get the first N most repeated characters ordered by the number of repetitions
func (x *CharStats) GetFirstN(n uint64) []uint8 {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	
	c := n
	if c > uint64(x.Count) { //Get the min betwen the given N and the characters count
		c = uint64(x.Count)
	}
	return x.Chars[0:c] //Return the first characters slice
}