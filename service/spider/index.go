package spider

import (
	jobInfo "internInformer/models/job_info"
)

type SpiderInterface interface {
	GetCompanyName() string
	UpdateJobInfo() (infos []jobInfo.JobInfo, err error)
	// 检查工作信息变化，如下架了某职位，新增了某职位
	CheckJobInfoChange() (newJobs, removedJobs []jobInfo.JobInfo, err error)
}

func GetDefaultSpider() SpiderInterface {
	return NewDefaultMihoyoInternSpider()
}
