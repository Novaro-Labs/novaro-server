package utils

import (
	"context"
	"github.com/robfig/cron/v3"
	"novaro-server/model"
)

type CronManager struct {
	cron *cron.Cron
}

func NewCronManager() *CronManager {
	return &CronManager{
		cron: cron.New(),
	}
}

func (cm *CronManager) AddJob(spec string, job func()) error {
	_, err := cm.cron.AddFunc(spec, job)
	return err
}

func (cm *CronManager) Start() {
	cli := model.GetRedisCli()
	_, err := cli.Ping(context.Background()).Result()
	if err == nil {
		cm.cron.Start()
	}
}

func (cm *CronManager) Stop() {
	cm.cron.Stop()
}
