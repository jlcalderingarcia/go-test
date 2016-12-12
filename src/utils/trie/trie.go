// Package trie implements a counting Trie to organize and count the words processed
// in the server implementation
package trie

import (
	"fmt"
)

// A Trie keeps the information of the text of the current node and the count of words which
// ends in the given node, the parent reference to form the words and the children nodes
// with their corresponding data
type Trie struct {
	Value string
	Count int64
	Parent *Trie
	Children *map[uint8]*Trie
}

// Initialize a new empty Trie node
func New() *Trie{
	res := new(Trie)
	res.Value = ""
	res.Count = 0
	res.Parent = nil
	res.Children = new(map[uint8]*Trie)
	*res.Children = make(map[uint8]*Trie)
	return res
}

// Conforms the word managed in the current node
func GetWord(t *Trie) string{
	if t.Parent == nil {
		return t.Value
	}
	return GetWord(t.Parent) + t.Value
}

// Insert a new word in the given Trie
func Insert(t *Trie, word string) *Trie{
	//if the word is the current one, just count it
	if t.Value == word { 
		t.Count++
		return t
	}
	
	index := findBreakIndex(t.Value, word) //find where we have to split the current text
	mymap := *t.Children 
	
	//Checks if the current node have to be broken
	if index < len(t.Value) && index > 0 {
		//Create the new node with the start part of the value
		newItem := New()
		newItem.Parent = t.Parent
		newItem.Value = t.Value[0:index]
		newItem.Count = 0
		if newItem.Parent != nil {
			parentmap := *newItem.Parent.Children
			parentmap[newItem.Value[0]] = newItem
		}
		//Update the current node
		t.Parent = newItem
		mymap1 := *newItem.Children
		mymap1[t.Value[index]] = t
		t.Value = t.Value[index:]
		
		//If the inserted word hasn't ended, insert the remaining text in a new node
		if index < len(word){
			otherItem := New()
			otherItem.Parent = newItem
			otherItem.Value = word[index:]
			otherItem.Count = 1
			mymap1[word[index]] = otherItem
			return otherItem
		} else {
			//If the word ends in the breaking text, update the counter
			newItem.Count = 1
			return newItem			
		}
	} else {
		rword := word[index:]
		child, exists := mymap[rword[0]]
		//Checks if there is a child with a common prefix with the word
		if exists {
			return Insert(child, rword)
		} else {
			//If no child with common preffix, insert a new child with the rest of the word
			newChild := New()
			newChild.Value = rword
			newChild.Count = 1
			newChild.Parent = t
			l := rword[0]
			mymap[l] = newChild
			return newChild
		}
	}
}

// Utility function to find the longest common preffix index
func findBreakIndex(t1, t2 string) int{
	mlen := len(t1)
	if len(t2) < mlen {
		mlen = len(t2)
	}
	i := 0
	for ; i < mlen; i++ {
		if t1[i] != t2[i] {
			return i
		}
	}
	return i
}

func Print(t *Trie){
	printTrie(t, "")
}
func printTrie(t *Trie, p string){
	fmt.Printf("%s%s - %d\n", p, t.Value, t.Count)
	for _, v := range *t.Children {
		printTrie(v, p + "-")
	}
}