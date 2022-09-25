package spider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMihoyo(t *testing.T) {
	spider := NewDefaultMihoyoInternSpider()
	infos, err := spider.UpdateJobInfo()
	assert.NoError(t, err)
	t.Log(infos)
}

func TestCheckMihoyoJobChange(t *testing.T) {
	spider := NewDefaultMihoyoInternSpider()
	newJobs, removedJobs, err := spider.CheckJobInfoChange()
	assert.NoError(t, err)
	t.Log(newJobs)
	t.Log(removedJobs)
}
