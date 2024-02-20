package set

type Set struct {
	m map[string]struct{}
}

func (s *Set) GetValues() []string {
	m := s.m
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func (s *Set) Add(str string) {
	s.m[str] = struct{}{}
}

func (s *Set) Remove(str string) {
	delete(s.m, str)
}

func (s *Set) Includes(str string) bool {
	_, included := s.m[str]
	return included
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]struct{})

	return s
}
