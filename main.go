package main

import (
	"SeedGo/logger"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/robfig/cron/v3"
	"io/ioutil"
)

type SeedConfig struct {
	Name        string   `json:"name"`
	LogLevel    string   `json:"log_level"`
	KafkaServer []string `json:"kafka_server"`
	TopicName   string   `json:"topic_name"`
	Schedule    []struct {
		Name         string                 `json:"name"`
		URL          string                 `json:"url"`
		Method       string                 `json:"method"`
		Headers      map[string]string      `json:"headers"`
		Body         map[string]interface{} `json:"body"`
		Timeout      int                    `json:"timeout"`
		Priority     int                    `json:"priority"`
		Retry        int                    `json:"retry"`
		Interval     int                    `json:"interval"`
		Tags         []string               `json:"tags"`
		BusinessType string                 `json:"business_type"`
	} `json:"schedule"`
}

type Seed struct {
	Name         string
	URL          string
	Method       string
	Headers      map[string]string
	Body         map[string]interface{}
	Timeout      int
	Priority     int
	Retry        int
	Interval     int
	Tags         []string
	BusinessType string
}

func main() {
	logger.Logger.Info("Starting...")

	// 创建调度器
	c := cron.New()
	// 添加种子生成任务
	spec := "* * * * *"
	_, err := c.AddFunc(spec, func() {
		logger.Logger.Info("任务开始")
		// 读取配置文件
		configFile, err := ioutil.ReadFile("config.json")
		if err != nil {
			logger.Logger.Error(err)
		}
		var seedConfig SeedConfig
		if err := json.Unmarshal(configFile, &seedConfig); err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to unmarshal configFile: %v", err))
		}

		// 构造种子列表
		var seeds []*Seed
		for _, item := range seedConfig.Schedule {

			seed := &Seed{
				Name:         item.Name,
				URL:          item.URL,
				Method:       item.Method,
				Headers:      item.Headers,
				Body:         item.Body,
				Timeout:      item.Timeout,
				Priority:     item.Priority,
				Retry:        item.Retry,
				Interval:     item.Interval,
				Tags:         item.Tags,
				BusinessType: item.BusinessType,
			}
			seeds = append(seeds, seed)
		}

		// 初始化 Kafka 生产者配置
		kafkaConfig := sarama.NewConfig()
		kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
		kafkaConfig.Producer.Retry.Max = 5
		kafkaConfig.Producer.Return.Successes = true

		// 连接 Kafka 集群
		kafkaServer := seedConfig.KafkaServer
		producer, err := sarama.NewSyncProducer(kafkaServer, kafkaConfig)
		if err != nil {
			logger.Logger.Error(fmt.Sprintf("Failed to create producer: %v", err))
		}
		defer func() {
			if err := producer.Close(); err != nil {
				logger.Logger.Error(fmt.Sprintf("Failed to close producer: %v", err))
			}
		}()

		topic := seedConfig.TopicName
		// 推送种子到kafka队列
		for _, seed := range seeds {
			data, err := json.Marshal(seed)
			if err != nil {
				logger.Logger.Error("Failed to marshal seed")
			}
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(data),
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("Failed to send message: %v", err))
			} else {
				logger.Logger.Info(fmt.Sprintf("topic: %v, partition: %v, offset: %v", topic, partition, offset))
			}
		}
	})
	if err != nil {
		panic(err)
		return
	}
	// 启动调度器
	c.Start()

	select {}
}
