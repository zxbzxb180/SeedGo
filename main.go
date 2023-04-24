package main

import (
	"SeedGo/handler"
	"SeedGo/logger"
	"github.com/robfig/cron/v3"
)

func main() {
	logger.Logger.Info("Starting...")

	// 创建调度器
	c := cron.New()
	// 添加种子生成任务
	spec := "* * * * *"
	_, err := c.AddFunc(spec, handler.SeedFromConfig)
	if err != nil {
		panic(err)
		return
	}
	// 启动调度器
	c.Start()

	select {}
}
