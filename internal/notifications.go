package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-msgsms-webhook/config"
	"huawei-msgsms-webhook/utils"
	"net/url"
	"strings"
)

// GrafanaAlert 结构体 grafana v6.4.3 告警json
type GrafanaAlert struct {
	EvalMatches []struct {
		Metric string   `json:"metric"`
		Tags   struct{} `json:"tags"`
		Value  int64    `json:"value"`
	} `json:"evalMatches"`
	ImageURL string `json:"imageUrl"`
	Message  string `json:"message"`
	RuleID   int64  `json:"ruleId"`
	RuleName string `json:"ruleName"`
	RuleURL  string `json:"ruleUrl"`
	State    string `json:"state"`
	Title    string `json:"title"`
}

// ProcessSMSVariables 处理用于短信模板中的变量
func ProcessSMSVariables(s string) string {
	// 如果为空，返回字符串 NULL
	// 如果不为空，按照指定字符串长度检查
	if utils.CheckStringUnicodeLength(s) == 0 {
		return "NULL"
	}
	return utils.TruncateStringByUnicodeLength(s, 20)
}

func EnvNameFromURL(rawURL string) string {
	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		slog.Error(fmt.Sprintln("URL解析失败:", err))
		return ""
	}

	// 从URL中获取路径部分
	path := parsedURL.Path

	// 将路径部分按照 "/" 进行分割
	segments := strings.Split(path, "/")

	// 如果路径包含至少4个部分，则取第4个部分作为结果
	// /d/xxxxx001-prod/xxxxx-prod 取最后的 xxxx-prod
	if len(segments) > 3 {
		result := segments[3]
		return result
	} else {
		slog.Error("URL不符合预期格式")
		return ""
	}
}

func alertNameFromRuleName(s string) string {
	//空格拆分字符串，取第二个
	// 使用空格分割字符串
	parts := strings.Split(s, " ")

	// 检查是否有足够的部分
	if len(parts) >= 2 {
		secondPart := parts[1]
		return secondPart
	} else {
		slog.Error("字符串格式不符合预期")
		return ""
	}
}

func envNameFromAlertName(s string) string {
	// 空格拆分字符串，取第1个
	// 使用空格分割字符串
	parts := strings.Split(s, " ")

	// 检查是否有足够的部分
	if len(parts) >= 2 {
		secondPart := parts[0]
		return secondPart
	} else {
		slog.Error("字符串格式不符合预期")
		return ""
	}
}

// AddDefaultCountryCode 检查电话号码是否包含国家码，如果不包含，则添加默认的国家码。
func AddDefaultCountryCode(phoneNumber string, defaultCountryCode string) string {
	// 假设国家码在电话号码中使用加号 "+" 表示
	if !strings.HasPrefix(phoneNumber, "+") {
		// 如果电话号码不以 "+" 开头，添加默认的国家码
		return defaultCountryCode + phoneNumber
	}
	// 否则，电话号码已包含国家码，不做修改
	return phoneNumber
}

func ProcessContactNumbers(sep string, values []string) string {
	defaultCountryCode := "+86" // 默认的国家码
	for i, v := range values {
		// 使用 AddDefaultCountryCode 函数来处理电话号码
		values[i] = AddDefaultCountryCode(v, defaultCountryCode)
	}
	return strings.Join(values, sep)
}

func SendNotifications(envName, alertName, status string) error {
	//	发送方式，可能存在sms或者别的，在这里统一处理

	slog.Info("提取内容",
		slog.String("envName", envName),
		slog.String("alertName", alertName),
		slog.String("state", status))

	envName = ProcessSMSVariables(envName)
	alertName = ProcessSMSVariables(alertName)
	status = ProcessSMSVariables(status)
	slog.Info("最终内容",
		slog.String("envName", envName),
		slog.String("alertName", alertName),
		slog.String("state", status))

	// 创建一个字符串切片
	strList := []string{envName, alertName, status}

	// 使用encoding/json包将切片转换为JSON格式的字符串
	jsonStr, err := json.Marshal(strList)
	if err != nil {
		slog.Error(fmt.Sprintln("JSON marshal error:", err))
		return err
	}
	// 输出JSON格式的字符串
	slog.Info(string(jsonStr))

	// 根据告警类别发送信息
	for _, alert := range config.Cfg.Alerts {
		if alert.Enabled == false {
			continue
		}
		switch t := alert.Type; t {
		case "huawei-msgsms":
			for _, receiver := range config.Cfg.Receivers {
				// 获取手机号
				if len(receiver.ContactNumbers) > 0 {
					receiverNumber := ProcessContactNumbers(",", receiver.ContactNumbers)
					slog.Info(fmt.Sprintf("发送告警信息到: %s", receiverNumber))
					if err := SendSMS(receiverNumber, string(jsonStr)); err != nil {
						slog.Error(fmt.Sprintln("Send SMS error:", err))
						return err
					}
				} else {
					slog.Info(fmt.Sprintf("接收方: %s,无联系方式,不发送告警信息.", receiver.Name))
				}
			}
		case "nmgykdxfsyy-sms":
			// todo: 内蒙古医院内部短信平台
		}
	}

	return nil
}

func HandleOldVersionNotification(alertData *map[string]interface{}) error {
	slog.Info("按照旧版消息结构处理")

	ruleUrl := (*alertData)["ruleUrl"].(string)
	ruleName := (*alertData)["ruleName"].(string)
	state := (*alertData)["state"].(string)

	envName := EnvNameFromURL(ruleUrl)
	alertName := alertNameFromRuleName(ruleName)

	// 发送消息
	if err := SendNotifications(envName, alertName, state); err != nil {
		return err
	}
	return nil
}

func HandleNewVersionNotification(alertData *map[string]interface{}) error {
	// 处理新版告警消息结构
	slog.Info("按照新版消息结构处理")

	alertsSlice, ok := (*alertData)["alerts"].([]interface{})
	if !ok {
		// 处理类型断言失败的情况
		return errors.New("无法将 alerts 转换为 []interface{}")
	}

	for i, v := range alertsSlice {
		alert, ok := v.(map[string]interface{})
		if !ok {
			// 处理类型断言失败的情况
			slog.Error(fmt.Sprintln("无法将 alert 转换为 map[string]interface{}"))
			continue
		}
		if i == 0 {
			status := alert["status"].(string)

			labels := alert["labels"].(map[string]interface{})
			alertName := labels["alertname"].(string)
			envName := labels["env_name"]

			if envName == nil {
				envName = envNameFromAlertName(alertName)
			}
			name := alertNameFromRuleName(alertName)
			// 发送消息
			if err := SendNotifications(envName.(string), name, status); err != nil {
				return err
			}
		}
	}

	return nil
}

func SelectSendNotifications(s []byte) error {
	// 判断警报信息结构属于老的还是新的，进行处理
	// alertData 动态结构体处理不同版本警报消息结构
	var alertData map[string]interface{}

	err := json.Unmarshal(s, &alertData)
	if err != nil {
		// 处理解析错误
		slog.Error(fmt.Sprintln("JSON marshal error:", err))
		return err
	}

	version := alertData["version"]

	if version == "1" {
		if err := HandleNewVersionNotification(&alertData); err != nil {
			slog.Error(fmt.Sprintln("HandleNewVersionNotification:", err))
			return err
		}
	}

	if version == nil {
		if err := HandleOldVersionNotification(&alertData); err != nil {
			slog.Error(fmt.Sprintln("HandleOldVersionNotification:", err))
			return err
		}
	}

	return nil
}
