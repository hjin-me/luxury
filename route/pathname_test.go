package route

import "testing"

func TestTrieCreate(t *testing.T) {
	tn := NewTrie()
	tn.Add([]string{"abc", "123", "xyz"})
	tn.Add([]string{"abc", "123", "tyu"})
	tn.Add([]string{"abc", "zxf", "tyu"})
	tn.Add([]string{"222", "zxf", "tyu"})
	if tn.String() != "[abc[123[xyz[]tyu[]]zxf[tyu[]]]222[zxf[tyu[]]]]" {
		t.Error(tn.String())
	}

}
