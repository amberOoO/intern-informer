package spider

import (
	"context"
	"database/sql"
	jobInfo "internInformer/models/job_info"
	"internInformer/utils"
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
	rawDB   *sql.DB
	jobQ    *jobInfo.Queries
}

func NewDefaultMihoyoInternSpider() *MihoyoInternSpider {
	client := &http.Client{}
	launcher := launcher.New().Headless(true).Devtools(true).NoSandbox(true)
	db := utils.GetSqliteDB()
	return &MihoyoInternSpider{
		logger:  utils.GetZapLogger(zap.String("service", "spider"), zap.String("spider", "mihoyo")),
		client:  client,
		browser: rod.New().ControlURL(launcher.MustLaunch()),
		rawDB:   db,
		jobQ:    jobInfo.New(db),
	}
}

func (s *MihoyoInternSpider) GetCompanyName() string {
	return "米哈游"
}

func (s *MihoyoInternSpider) UpdateJobInfo() (infos []jobInfo.JobInfo, err error) {
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

		infos = append(infos, jobInfo.JobInfo{
			JobName:     jobName,
			JobType:     sql.NullString{String: jobType, Valid: true},
			JobCategory: sql.NullString{String: jobCategory, Valid: true},
			JobLocation: sql.NullString{String: location, Valid: true},
		})
	}

	// 首先将is_exists字段清空，为后续更新做准备
	err = s.jobQ.ResetIsExistsByCompany(context.Background(), s.GetCompanyName())
	if err != nil {
		s.logger.Error("error reset is_exists", zap.Error(err))
		return nil, err
	}

	// 上传数据库咯
	for _, info := range infos {
		// 如果数据库中已经存在该职位，则更新is_exists字段。否则插入
		err = s.jobQ.InsertJobInfo(context.Background(), jobInfo.InsertJobInfoParams{
			JobName:     info.JobName,
			Company:     s.GetCompanyName(),
			JobType:     info.JobType,
			JobCategory: info.JobCategory,
			JobLocation: info.JobLocation,
		})
		if err != nil {
			s.logger.Info("InsertJobInfo failed", zap.Error(err))
		}
	}

	return infos, nil
}

func (s *MihoyoInternSpider) CheckJobInfoChange() (newJobs, removedJobs []jobInfo.JobInfo, err error) {
	_, _ = s.UpdateJobInfo()
	tx, err := s.rawDB.Begin()
	if err != nil {
		s.logger.Error("error begin transaction", zap.Error(err))
		return nil, nil, err
	}
	defer tx.Rollback()
	qtx := s.jobQ.WithTx(tx)
	/* 获取新增加（未checked的job） */
	newJobs, err = qtx.GetJobInfosByIsCheckedAndCompany(context.Background(), jobInfo.GetJobInfosByIsCheckedAndCompanyParams{
		Company:   s.GetCompanyName(),
		IsChecked: false,
	})
	if err != nil {
		s.logger.Error("get new job info from db failed", zap.Error(err))
		return nil, nil, err
	}

	/* 获取已经被删除的数据（is_exists为false） */
	removedJobs, err = qtx.GetJobInfosByIsExistsAndCompany(context.Background(), jobInfo.GetJobInfosByIsExistsAndCompanyParams{
		Company:  s.GetCompanyName(),
		IsExists: false,
	})
	if err != nil {
		s.logger.Error(
			"get removed job info from db failed",
			zap.String("company", s.GetCompanyName()),
			zap.Error(err),
		)
		return newJobs, removedJobs, err
	}

	/* 更新状态为已经checked */
	for _, job := range newJobs {
		err := qtx.UpdateIsChecked(context.Background(), jobInfo.UpdateIsCheckedParams{
			ID:        job.ID,
			IsChecked: true,
		})
		if err != nil {
			s.logger.Info(
				"update is_checked failed",
				zap.String("job_name", job.JobName),
				zap.String("company", job.Company),
				zap.Error(err),
			)
		}
	}

	err = tx.Commit()
	if err != nil {
		s.logger.Error("error commit transaction", zap.Error(err))
		return nil, nil, err
	}
	return newJobs, removedJobs, nil
}
