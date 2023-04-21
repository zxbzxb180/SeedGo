package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Name     string `json:"name"`
	LogLevel string `json:"log_level"`
	Schedule []struct {
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
	log.Println("Starting...")

	// 创建调度器
	c := cron.New()
	// 添加种子生成任务
	spec := "* * * * *"
	_, err := c.AddFunc(spec, func() {
		log.Println("任务开始")

		// 初始化 Kafka 生产者配置
		producerConfig := &kafka.ConfigMap{
			"bootstrap.servers":     "localhost:9092", // Kafka 服务器地址
			"retries":               3,                // 自动重试次数
			"retry.backoff.ms":      100,              // 重试时间间隔
			"session.timeout.ms":    10000,            // 心跳超时时间
			"heartbeat.interval.ms": 5000,             // 发送心跳消息的时间间隔
		}

		// 创建 Kafka 生产者
		producer, err := kafka.NewProducer(producerConfig)
		if err != nil {
			fmt.Printf("Failed to create Kafka producer: %s\n", err)
			return
		}
		defer producer.Close()

		// 读取配置文件
		configFile, err := ioutil.ReadFile("config.json")
		if err != nil {
			panic(err)
		}
		var config Config
		if err := json.Unmarshal(configFile, &config); err != nil {
			panic(err)
		}

		// 构造种子列表
		var seeds []*Seed
		for _, item := range config.Schedule {

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

		// 输出种子列表
		for _, seed := range seeds {
			log.Printf("%+v\n", seed)
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
