package scheduler

import (
	"bytes"
	jobInfo "internInformer/models/job_info"
	"internInformer/utils"
	"time"

	"text/template"

	"go.uber.org/zap"
)

type DefaultScheduler struct {
	logger    *zap.Logger
	instances []*Instance
	timer     *time.Timer
}

func NewDefaultScheduler() *DefaultScheduler {
	vi := utils.GetDefaultViper()
	schedCfg := vi.GetStringMap("scheduler")

	return &DefaultScheduler{
		timer: time.NewTimer(time.Duration(schedCfg["check_interval"].(int64)) * time.Minute),
		// timer:  time.NewTimer(time.Duration(5 * time.Second)),
		logger: utils.GetZapLogger(),
	}
}

func (ds *DefaultScheduler) Register(ins *Instance) {
	ds.instances = append(ds.instances, ins)
}

func (ds *DefaultScheduler) Start() {
	select {
	case <-ds.timer.C:
		ds.logger.Info(
			"start to schedule",
			zap.String("time", time.Now().String()),
		)
		for _, ins := range ds.instances {
			go ds.runInstance(ins)
		}
	}
}

func (ds *DefaultScheduler) runInstance(ins *Instance) {
	ds.logger.Info(
		"start check info change",
		zap.String("spider", ins.Spider.GetCompanyName()),
	)
	newJobs, removedJobs, err := ins.Spider.CheckJobInfoChange()
	if err != nil {
		ds.logger.Error(
			"failed to check job info change",
			zap.String("spiderCompany", ins.Spider.GetCompanyName()),
			zap.Error(err),
		)
	}
	if len(newJobs) == 0 && len(removedJobs) == 0 {
		return
	}
	summary := ds.generateJobInfoChangesSummary(newJobs, removedJobs)
	err = ins.Informer.Send(summary)
	if err != nil {
		ds.logger.Error(
			"failed to send message",
			zap.String("informer", ins.Spider.GetCompanyName()),
			zap.Error(err),
		)
	}
}

var summaryTmpl = `
# 工作变动如下
## 1. 新增
{{ range $k, $v := .NewJobs }}
{{$k}}. {{ $v.Company }} - {{ $v.JobName }} - {{ $v.JobType }} - {{ $v.JobType }} - {{ $v.JobCategory }} - {{ $v.JobLocation }}
{{ else }}
没有新增的工作
{{ end }}

## 2. 移除
{{ range $k, $v := .RemovedJobs }}
{{$k}}. {{ $v.Company }} - {{ $v.JobName }} - {{ $v.JobType }} - {{ $v.JobType }} - {{ $v.JobCategory }} - {{ $v.JobLocation }}
{{ else }}
没有移除的工作
{{ end }}
`

func (ds *DefaultScheduler) generateJobInfoChangesSummary(newJobs, removedJobs []jobInfo.JobInfo) string {
	var (
		resultBuffer *bytes.Buffer = bytes.NewBuffer([]byte{})
	)
	tmpl, err := template.New("job info changes summary").Parse(summaryTmpl)
	if err != nil {
		ds.logger.Error("failed to parse template", zap.Error(err))
		return ""
	}
	err = tmpl.Execute(resultBuffer, map[string][]jobInfo.JobInfo{
		"NewJobs":     newJobs,
		"RemovedJobs": removedJobs,
	})
	if err != nil {
		ds.logger.Error("failed to execute template", zap.Error(err))
		return ""
	}
	return resultBuffer.String()
}
