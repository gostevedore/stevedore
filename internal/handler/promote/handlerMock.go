package handler

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type HandlerMock struct {
	mock.Mock
}

func NewHandlerMock() *HandlerMock {
	return &HandlerMock{}
}

func (h *HandlerMock) Handler(ctx context.Context, options *Options) error {
	args := h.Called(ctx, options)

	return args.Error(0)
}
