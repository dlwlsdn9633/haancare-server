package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type CronJob struct {
	Name     string // 작업 이름
	Schedule string // 크론 주기
	Task     func() // 실제 실행할 함수
}

func InitCronJobs() (err error) {
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	jobList := []CronJob{
		{
			Name: "SetInvoiceNumber",
			//Schedule: "0 */2 * * * *",
			Schedule: "0 */2 * * * *",
			Task:     CronSetInvoiceNumber,
		},
		{
			Name:     "SetSessionID",
			Schedule: "0 */15 * * * *",
			Task:     CronSetSessionID,
		},
		{
			Name:     "SendNateOnMsg",
			Schedule: "0 0 */3 * * *",
			Task:     CronSendNateOnMsg,
		},
	}

	for _, job := range jobList {
		_, err = c.AddFunc(job.Schedule, func() {
			logger.Info("Cron Task Started", "jobName", job.Name)
			job.Task()
		})
		if err != nil {
			logger.Error("Can Registration Failed", "jobName", job.Name)
			err = errors.Wrap(err, "cron registration failed")
			return
		}
	}

	c.Start()
	logger.Info("Cron Start")
	return
}

func CronSetInvoiceNumber() {
	ordersNum, err := GetOrderNums(OrderStatPrint)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get ordersNum: %+v", err))
		return
	}

	baseUrls := []string{
		"https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrgstinfo",
		"https://pid.alps.llogis.com:18210/pid/ftr/pacltrc/inner/popinvrsrvinfo",
	}

	for _, orderNum := range ordersNum {
		var allResults []OrderResult
		for _, url := range baseUrls {
			results, err := GetAlpsOrders(url, orderNum, alpsToken)
			if err != nil {
				logger.Error(fmt.Sprintf("API call failed: %v (url: %s, orderNum: %s)", err, url, orderNum))
				continue
			}
			allResults = append(allResults, results...)
		}

		if len(allResults) == 0 {
			// TODO: 네이트온 메시지 보내기
			logger.Error(fmt.Sprintf("empty order results after checking all URLs (orderNum: %s)", orderNum))
			continue
		}

		firstResult := allResults[0]
		cnt, err := UpdateOrderInvoice(orderNum, firstResult.InvNo)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update delivery info: %+v (orderNum: %s)", err, orderNum))
			continue
		}

		if cnt == 0 {
			logger.Warn(fmt.Sprintf("No order found to update: %s", orderNum))
			continue
		}

		logger.Info(fmt.Sprintf("Success update: %s -> %s", orderNum, firstResult.InvNo))
	}
}

func CronSetSessionID() {
	var err error
	alpsToken, err = GetLatestAlpsToken()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to refresh session id: %+v", err))
		return
	}
}

// TODO: nateon 알림으로 특정 시간동안 어떤 것이 연동되지 않았는지 톡보내는 Cron 작업 추가 (작업 보고용으로 남기면 좋을듯?)
func CronSendNateOnMsg() {

}
