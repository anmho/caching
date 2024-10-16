package todo

import "github.com/anmho/caching/cache"

type Service struct {
	cacheStrategy cache.Strategy
}

func WithCacheStrategy(strategy cache.Strategy) func(s *Service) {
	return func(s *Service) {
		s.cacheStrategy = strategy
	}
}

func MakeService(opts ...func(o *Service)) *Service {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Service) CreateTodo() {

}

func (s *Service) GetTodo() {

}

func (s *Service) UpdateTodo() {

}

func (s *Service) DeleteTodo() {

}
