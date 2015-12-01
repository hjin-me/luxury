package route

import "log"

type trieNode struct {
	Frag  string
	Child []*trieNode
}

func NewTrie() trieNode {
	return trieNode{}
}

func (tn *trieNode) Add(frags []string) {
	if len(frags) == 0 {
		return
	}

	var ntn *trieNode
	for _, t := range tn.Child {
		log.Println(t.Frag, frags[0])
		if t.Frag == frags[0] {
			ntn = t
		}
	}
	if ntn == nil {
		log.Println("ntn is nil")
		ntn = &trieNode{}
		ntn.Frag = frags[0]
		ntn.Child = make([]*trieNode, 0)
		tn.Child = append(tn.Child, ntn)
	}
	if len(frags) > 1 {
		ntn.Add(frags[1:])
	}
}

func (tn trieNode) String() string {
	s := tn.Frag + "["
	for _, t := range tn.Child {
		s = s + t.String()
	}
	return s + "]"
}
