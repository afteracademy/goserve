package network

import (
	"context"
)

type baseService struct {
	context context.Context
}

func NewBaseService() BaseService {
	return &baseService{
		context: context.Background(),
	}
}

func (s *baseService) Context() context.Context {
	return s.context
}
