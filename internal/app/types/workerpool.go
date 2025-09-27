package types

type WorkerPool interface {
	Submit(func()) error
}
