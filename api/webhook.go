package api

import (
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-msgsms-webhook/config"
	"huawei-msgsms-webhook/internal"
	"io"
	"net/http"
)

func getServerPort() int {
	var port int

	cfg := config.Cfg

	if cfg.Server.Port == 0 {
		port = 8080
	} else {
		port = cfg.Server.Port
	}
	return port
}

func StartWebhookServer() error {
	port := getServerPort()
	http.HandleFunc("/webhook", handleWebhook)
	//port := 8080
	slog.Info(fmt.Sprintf("Webhook server listening on port %d ...", port))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		return err
	}
	return nil
}

func EchoConfig() {
	cfg := config.Cfg
	// 使用配置中的值
	var port int
	if cfg.Server.Port == 0 {
		port = 8080
	} else {
		port = cfg.Server.Port
	}

	fmt.Printf("Server Port: %d\n", port)
	for _, v := range cfg.Alerts {
		fmt.Printf("%s\n", v.Type)
		fmt.Printf("%t\n", v.Enabled)
	}
	for _, v := range cfg.Receivers {
		fmt.Printf("%s\n", v.Name)
		fmt.Printf("len: %d", len(v.ContactNumbers))
		for _, num := range v.ContactNumbers {
			fmt.Printf("%s\n", num)
		}
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	slog.Info(fmt.Sprintf("Received Webhook Payload:\n%s\n", body))

	// 在这里处理来自 Grafana 的告警数据
	// 您可以解析 JSON 数据，触发相应的操作

	if err := internal.SelectSendNotifications(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "通知发送失败")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received and processed successfully")
}
