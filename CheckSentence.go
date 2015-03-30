package spellcheck

import (
	"fmt"
)

type SentenceCorrection struct {
	BeginAt int
	EndAt   int
	Result  *CheckResult
}

func (sc *SentenceCorrection) Apply(s string) string {
	fmt.Printf("SC: %#v Word: \"%s\"\n", sc, sc.Result.Correction)
	//fmt.Printf("Removing: %s\n", s[sc.BeginAt:sc.EndAt]);
	suffix := ""
	if len(s) > sc.EndAt {
		suffix = s[sc.EndAt:]
	}
	return s[:sc.BeginAt] + sc.Result.Correction + suffix
}

type SentenceCorrectionSet []SentenceCorrection

func (scs SentenceCorrectionSet) ApplyAll(s string) string {
	for _, sc := range scs {
		println(s)
		s = sc.Apply(s)
	}
	return s
}

func (d *Dict) CheckSentence(sentence string) SentenceCorrectionSet {
	res := []SentenceCorrection{}
	beganat := 0
	for i := 0; i <= len(sentence); i++ {
		if i < len(sentence) && sentence[i] != ' ' {
			continue
		}

		word := sentence[beganat:i]
		r := d.CheckWord(word)
		if r.IsCorrect == false && len(r.Correction) > 0 {
			res = append(res, SentenceCorrection{
				BeginAt: beganat + 1,
				EndAt:   i,
				Result:  r,
			})
		}

		beganat = i
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
