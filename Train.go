package spellcheck

import (
	"bufio"
	"io"
	"os"
	"strings"
)

var DefaultDelimiters = " _=.?!,\n\"'"

// TrainCmd is a command ran within a dictionary's goroutine that
// populates the internal dictionary with words.
type TrainCmd struct {
	Reader     io.Reader
	Delimiters string
    Done chan bool
}

// Run() is the internal TrainCmd function to fulfil the DictCmd
// interface. This is executed within the Dict worker goroutine, and
// should not be called directly.
func (cmd *TrainCmd) Run(d *Dict) {
	defer close(cmd.Done);
    r := bufio.NewReader(cmd.Reader)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}
		line = strings.ToLower(line)
		words := strings.Split(line, " ")
		for _, word := range words {
			word = strings.Trim(word, cmd.Delimiters)
			if _, found := d.Dict[word]; !found{
				d.Dict[word] = 1
			} else {
				d.Dict[word]++
			}
		}
	}
}

// TrainFile() accepts a path to a file, opens that file, and will create
// a TrainCmd and queue it.
func (d *Dict) TrainFile(filename string) (error, chan bool) {
	file, e := os.Open(filename)
	if e != nil {
		return e, nil
	}
	return nil, d.TrainReader(file)
}

// TrainReader() is identical to TrainDict() except a io.Reader is
// provided instead of a file path.
func (d *Dict) TrainReader(r io.Reader) chan bool {
    c := make(chan bool);
    d.PushQueue(&TrainCmd{
		Reader:     r,
		Delimiters: DefaultDelimiters,
        Done: c,
	})
    return c;
}
