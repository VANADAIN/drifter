package network

type Filter struct {
	Blacklist []string // addresses
	Tags      []string
	Banwords  []string
}

func NewFilter(blacklist []string, tags []string, badwords []string) {
	//
}

func (f *Filter) FilterMessage() {
	//
}
