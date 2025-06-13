package main

type ResolverStack []map[string]bool

func (s ResolverStack) IsEmpty() bool {
	if len(s) == 0 {
		return true
	}
	return false
}

func (s *ResolverStack) Push(m map[string]bool) {
	*s = append(*s, m)
}

func (s *ResolverStack) Pop() (map[string]bool, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	index := len(*s) - 1
	elem := (*s)[index]
	*s = (*s)[:index]
	return elem, true
}

func (s *ResolverStack) Peek() (*map[string]bool, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	return &(*s)[len(*s)-1], true
}
