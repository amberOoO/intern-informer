package informer

type informerInterface interface {
	Send(msg string) error
}
