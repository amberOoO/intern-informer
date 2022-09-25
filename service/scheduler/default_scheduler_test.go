package scheduler

import (
	jobInfo "internInformer/models/job_info"
	"internInformer/service/informer"
	"internInformer/service/spider"
	"testing"
)

func TestDefaultScheduler(t *testing.T) {
	s := NewDefaultScheduler()
	s.Register(&Instance{
		Spider:   spider.NewDefaultSpider(),
		Informer: informer.NewDefaultInformer(),
	})
	s.Start()
}

func TestGenerateJobInfoChangesSummary(t *testing.T) {
	mockJobInfos := []jobInfo.JobInfo{
		{
			JobName: "测试",
			Company: "米哈游",
		},
	}
	s := NewDefaultScheduler()
	summary := s.generateJobInfoChangesSummary(mockJobInfos, mockJobInfos)
	t.Log(summary)
}
