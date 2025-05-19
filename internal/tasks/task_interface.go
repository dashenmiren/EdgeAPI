package tasks

type TaskInterface interface {
	Start() error
	Loop() error
	Stop() error
}
