package store

import "sync"

type RuntimeStorege struct {
	mx sync.Mutex
	m  map[string][]string
}

func NewRuntimeStorege() *RuntimeStorege {
	return &RuntimeStorege{
		m:  make(map[string][]string),
		mx: sync.Mutex{},
	}
}

func (s *RuntimeStorege) Load(key string) ([]string, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	val, ok := s.m[key]
	return val, ok
}

func (s *RuntimeStorege) Pop(key string) (string, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	vals, ok := s.m[key]
	var x string
	if ok && (len(vals) > 1) {
		x, vals = vals[0], vals[1:]
		s.m[key] = vals
	} else if len(vals) == 1 {
		x = vals[0]
		delete(s.m, key)
	}
	return x, ok
}

func (s *RuntimeStorege) Exist(key string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()

	val, ok := s.m[key]
	return ok && len(val) > 0
}

func (s *RuntimeStorege) Store(key string, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if val, ok := s.m[key]; ok != true {
		s.m[key] = []string{value}
	} else {
		s.m[key] = append(val, value)
	}
}
