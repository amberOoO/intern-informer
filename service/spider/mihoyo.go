package spider

import (
	"internInformer/service/utils"
	"net/http"

	"github.com/go-rod/rod"

	"github.com/go-rod/rod/lib/launcher"
	"go.uber.org/zap"
)

var MihoyoInternPageUrl = "https://campus.mihoyo.com/?recommendationCode=4XKW&isRecommendation=true#/campus/position?jobNatures%5B0%5D=3"

type MihoyoInternSpider struct {
	logger  *zap.Logger
	browser *rod.Browser
	client  *http.Client
}

func NewDefaultMihoyoInternSpider() *MihoyoInternSpider {
	client := &http.Client{}
	launcher := launcher.New().Headless(true).Devtools(true).NoSandbox(true)
	return &MihoyoInternSpider{
		logger:  utils.GetZapLogger(zap.String("service", "spider"), zap.String("spider", "mihoyo")),
		client:  client,
		browser: rod.New().ControlURL(launcher.MustLaunch()),
	}
}

func (s *MihoyoInternSpider) GetJobInfo() (infos []JobInfo, err error) {
	var (
		browserInst = s.browser.MustConnect()
	)
	defer browserInst.Close()

	page := browserInst.MustPage(MihoyoInternPageUrl).MustWaitLoad()

	// Start to analyze request events
	wait := page.MustWaitRequestIdle()
	wait()

	elements := page.MustElements(".jobItem")
	// 获取所有实习列表中的item
	for _, elem := range elements {
		// 解析职位名
		jobName := elem.MustElement(".jobName").MustText()
		// 解析详情
		details := elem.MustElements(".contentText")
		location := details[0].MustText()
		jobCategory := details[1].MustText()
		jobType := details[2].MustText()

		infos = append(infos, JobInfo{
			JobName:     jobName,
			JobType:     jobType,
			JobCategory: jobCategory,
			Location:    location,
		})
	}
	return infos, nil
}
