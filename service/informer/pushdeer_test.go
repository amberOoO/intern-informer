package informer

import "testing"

func TestPushdeerSend(t *testing.T) {
	pi := NewDefaultPushdeer()
	pi.Send("# 测试\n## 新增职位有\n - 1\n - 2\n## 下架职位有\n - 3\n - 4")
}
