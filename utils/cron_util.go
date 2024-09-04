package utils

import (
	"github.com/robfig/cron/v3"
	"novaro-server/config"
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

	rdb := config.RDB
	_, err := rdb.Ping(rdb.Context()).Result()
	if err == nil {
		cm.cron.Start()
	}
}

func (cm *CronManager) Stop() {
	cm.cron.Stop()
}
