package trie

import (
	"testing"
)

// Test the Init function ensuring the initialized node has empty value, counter 0, ni parent and empty children map
func TestInit(t *testing.T){
	tr := New()
	if tr.Value != "" {
		t.Error("Error, the value wasn't initialize empty")
	}
	if tr.Count != 0 {
		t.Error("Error, the count wasn't initialize on 0")
	}
	if tr.Parent != nil {
		t.Error("Error, the Parent wasn't initialize nil")
	}
	if tr.Children == nil {
		t.Error("Error, the Children was initialize nil")
	}
	if len(*tr.Children) > 0 {
		t.Error("Error, the Children wasn't initialize empty")
	}
}

//To validate the GetWord function works propertly
func TestGetWord(t *testing.T){
	//Create a basic 2 nodes Trie
	tr := New()
	tr.Value = "amaz"
	end := New()
	end.Value = "ing"
	end.Parent = tr
	children := *tr.Children
	children[end.Value[0]] = end
	
	if GetWord(end) != "amazing" {
		t.Errorf("Error, se obtuvo %s cuando se esperaba %s", GetWord(end), "amazing")
	}
}

//Insert several words, including breaking, ending, different starting and all cases
//and validate the inserted word is ok, as well as the counters
func TestInsert(t *testing.T){
	//Test cases values
	words := []string{ "test", "testing", "tester", "hello", "go", "world", "tellerman", "tellerma", "test", "hell", "test" }
	counts := []int64 { 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 3 }
	
	//Initialize an empty Trie
	tr := New()
	
	//Insert the test cases values
	for i, w := range words {
		t1 := Insert(tr, w)
		if GetWord(t1) != w {
			t.Errorf("Error, se obtuvo %s cuando se esperaba %s", GetWord(t1), w)
		}
		if t1.Count != counts[i] {
			t.Errorf("Error, %s se ha contado %d veces en vez de %d veces", w, t1.Count, counts[i])
		}
	}
}