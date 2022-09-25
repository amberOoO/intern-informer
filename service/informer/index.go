package informer

type InformerInterface interface {
	Send(msg string) error
}

func NewDefaultInformer() InformerInterface {
	return NewDefaultPushdeer()
}
