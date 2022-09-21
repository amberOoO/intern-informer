package spider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMihoyo(t *testing.T) {
	spider := NewDefaultMihoyoInternSpider()
	infos, err := spider.GetJobInfo()
	assert.NoError(t, err)
	t.Log(infos)
}
