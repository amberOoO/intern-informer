package spider

import (
	jobInfo "internInformer/models/job_info"
)

type PageInfo struct {
	Url      string
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// 工作类型，实习或全职
type SpiderInterface interface {
	GetCompanyName() string
	UpdateJobInfo() (infos []jobInfo.JobInfo, err error)
	// 检查工作信息变化，如下架了某职位，新增了某职位
	CheckJobInfoChange() (newJobs, removedJobs []jobInfo.JobInfo, err error)
}
