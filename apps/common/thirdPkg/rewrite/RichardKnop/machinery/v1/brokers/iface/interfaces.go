package iface

import (
	"context"

	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/config"
	"WenBeego/apps/common/thirdPkg/rewrite/RichardKnop/machinery/v1/tasks"
)

type DLQHandler func(signature *tasks.Signature) error

// Broker - a common interface for all brokers
type Broker interface {
	GetConfig() *config.Config
	SetRegisteredTaskNames(names []string)
	IsTaskRegistered(name string) bool
	StartConsuming(consumerTag string, concurrency int, p TaskProcessor) (bool, error)
	StopConsuming()
	Publish(ctx context.Context, task *tasks.Signature) error
	GetPendingTasks(queue string) ([]*tasks.Signature, error)
	GetDelayedTasks() ([]*tasks.Signature, error)
	AdjustRoutingKey(s *tasks.Signature)
	StartDLQConsuming(dlqQueue string, handler DLQHandler) error
}

// TaskProcessor - can process a delivered task
// This will probably always be a worker instance
type TaskProcessor interface {
	Process(signature *tasks.Signature) error
	CustomQueue() string
	PreConsumeHandler() bool
}
