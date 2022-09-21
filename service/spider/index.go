package spider

type PageInfo struct {
	Url      string
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// 工作类型，实习或全职
type JobType int

type JobTime int

type JobInfo struct {
	JobName     string
	Description string
	// 实习或者全职
	JobType string
	// 工作分类，如计算机、综合等
	JobCategory string
	// 工作地点
	Location string
}

type SpiderInterface interface {
	GetJobInfo() (infos []JobInfo, err error)
}
