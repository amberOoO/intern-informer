package informer

import "testing"

func TestPushdeerSend(t *testing.T) {
	pi := NewDefaultPushdeer()
	pi.Send("go test message")
}
