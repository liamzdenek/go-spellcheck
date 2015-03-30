package spellcheck

// Dict is the central struct for this library. It maintains the worker
// goroutines, and contains the internal data. All member functions of
// dict are thread safe.
type Dict struct {
	Cmds chan DictCmd
	Dict map[string]int
}

// NewDict() will create and initialize a new Dictionary, and start the
// worker goroutine
func NewDict() *Dict {
	d := &Dict{
		Cmds: make(chan DictCmd),
		Dict: map[string]int{},
	}
	go func() {
		for cmd := range d.Cmds {
			cmd.Run(d)
		}
	}()
	return d
}

// Known() will return an integer value for whether the word exists in
// the dictionary. 0 means that the word was not found.
func (d *Dict) Known(word string) int {
	val, found := d.Dict[word]
	if !found {
		return 1 // treat novel words as if we've seen them once
	}
	return val
}

// PushQueue() will push a DictCmd task into the command queue. This is
// not a guarantee that it will be executed now or next, but that it
// will be executed eventually
func (d *Dict) PushQueue(cmd DictCmd) {
	d.Cmds <- cmd
}

// Destroy() will terminate all worker goroutines. Once you release all
// references to the object, the internal data will clean up, too.
func (d *Dict) Destroy() {
	close(d.Cmds)
}

// DictCmd is an interface that all internal commands intending on
// modifying internal Dict data must fulfill. This is to ensure
// multithread safety
type DictCmd interface {
	Run(d *Dict)
}
