package messagebroker

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Publish(ctx context.Context, topic string, message interface{}) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}