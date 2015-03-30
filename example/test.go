package main

import (
	".."
	"fmt"
)

func main() {
	dict := spellcheck.NewDict()
	Check(dict.TrainFile("big.txt"))
	//    loe = love   -- char deletion
	//    ypu = you    -- char replacement
	// jeesus = jesus  -- char insertion
	// chrsit = christ -- char transposition
	sentence := "i loe ypu jeesus chrsit"
	corrections := dict.CheckSentence(sentence)

	fmt.Printf("Corrected: %s\n", corrections.ApplyAll(sentence))
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
