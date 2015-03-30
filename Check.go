package spellcheck

import (
	"strings"
)

var DefaultAlphabet = strings.Split("abcdefghijklmnopqrstuvwxyz", "")

type CheckResult struct {
	IsCorrect  bool
	Correction string
	Score      int
}

type CheckCmd struct {
	Word     string
	Alphabet []string
	Result   chan *CheckResult
}

// Run() is the logic that will be executed by the Dict worker thread.
// This should not be called directly. It will determine if the provided
// word is in the dictionary, and, if it isn't, it will generate a set of
// variants using Variants(), and will send the best result down
// cmd.Result in the form of a CheckResult
func (cmd *CheckCmd) Run(d *Dict) {
	if score, found := d.Dict[cmd.Word]; found && score > 0 {
		cmd.Result <- &CheckResult{
			IsCorrect:  true,
			Correction: cmd.Word,
			Score:      score,
		}
		return
	}

	variants := cmd.Variants(cmd.Word)

	best, score := cmd.GetBestMatch(d, variants)

	cmd.Result <- &CheckResult{
		IsCorrect:  false,
		Correction: best,
		Score:      score,
	}
}

// Variants() will take a single word and generate a set of possible
// permutations
func (cmd *CheckCmd) Variants(word string) []string {
	out := make(chan string)

	go func() {
		defer close(out)

		// definitions
		type Split struct {
			A string
			B string
		}

		// splits
		splits := []Split{}
		for i := 0; i < len(word)+1; i++ {
			splits = append(splits, Split{
				A: strings.Trim(word[:i], " "),
				B: strings.Trim(word[i:], " "),
			})
		}

		for _, set := range splits {
			// deletes
			if len(set.B) >= 1 {
				out <- set.A + set.B[1:]
			}
			// transposes
			if len(set.B) >= 2 {
				out <- set.A + string(set.B[1]) + string(set.B[0]) + set.B[2:]
			}

			for _, letter := range cmd.Alphabet {
				// replaces
				if len(set.B) >= 1 {
					out <- set.A + letter + set.B[1:]
				}
				// inserts
				out <- set.A + letter + set.B
			}
		}
	}()

	variants := []string{}
	for word := range out {
		variants = append(variants, word)
	}
	return variants
}

// GetBestMatch() will take a set of variants and return the one with the
// highest score per the dict.Known() function
func (cmd *CheckCmd) GetBestMatch(d *Dict, variants []string) (best_string string, best_score int) {
	for _, variant := range variants {
		score := d.Known(variant)
		if score > best_score {
			best_score = score
			best_string = variant
		}
	}
	return best_string, best_score
}

// CheckWord() is a primitive used to determine if a single word is
// spelled correctly
func (d *Dict) CheckWord(word string) *CheckResult {
	cmd := &CheckCmd{
		Word:     word,
		Alphabet: DefaultAlphabet,
		Result:   make(chan *CheckResult),
	}
	d.PushQueue(cmd)
	return <-cmd.Result
}
