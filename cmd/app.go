package main

import (
	"flag"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-msgsms-webhook/api"
	"huawei-msgsms-webhook/config"
	"log"
)

var configFile string

func init() {
	// 使用 flag 包定义 -c 参数，指定配置文件路径
	flag.StringVar(&configFile, "c", "/opt/webhook/config.yaml", "Path to the configuration file")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 打印配置文件路径，供你后续使用
	slog.Info(fmt.Sprintf("Using config file: %s", configFile))

	// 初始化配置
	if err := config.InitConfig(configFile); err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// 启动 Webhook 服务
	//api.EchoConfig()
	err := api.StartWebhookServer()
	if err != nil {
		slog.Error(fmt.Sprintf("Error starting server: %v", err))
	}
}
