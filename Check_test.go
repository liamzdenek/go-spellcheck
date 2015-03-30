package spellcheck

import (
	"testing"
)

func BenchmarkSingleWord(b *testing.B) {
	dict := NewDict()
	e, c := dict.TrainFile("example/big.txt")
	Check(e)
	<-c
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.CheckWord("chrsit")
	}
}

func BenchmarkSentence(b *testing.B) {
	dict := NewDict()
	e, c := dict.TrainFile("example/big.txt")
	Check(e)
	<-c
	sentence := "i loe ypu jeesus chrsit"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.CheckSentence(sentence)
	}
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
