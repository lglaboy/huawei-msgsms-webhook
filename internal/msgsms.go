package internal

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/exp/slog"
	"huawei-msgsms-webhook/config"
	core "huaweicloud.com/apig/signer"
	"io"
	"net/http"
	"net/url"
)

// SendSMS 发送SMS信息
func SendSMS(receiver, templateParas string) error {
	hw := config.Cfg.Huawei

	//必填,请参考"开发准备"获取如下数据,替换为实际值
	appInfo := core.Signer{
		Key:    hw.AppKey,    //App Key
		Secret: hw.AppSecret, //App Secret
	}

	body := buildRequestBody(
		//国内短信签名通道号
		hw.Sender,
		//短信接收人
		receiver,
		//模板ID
		hw.TemplateId,
		templateParas,
		//选填,短信状态报告接收地址,推荐使用域名,为空或者不填表示不接收状态报告
		hw.StatusCallBack,
		//签名名称
		hw.Signature,
	)

	// 发送
	resp, err := post(hw.ApiAddress, []byte(body), appInfo)

	if err != nil {
		return err
	}
	slog.Info(resp)

	return nil
}

/**
 * sender,receiver,templateId不能为空
 */
func buildRequestBody(sender, receiver, templateId, templateParas, statusCallBack, signature string) string {
	param := "from=" + url.QueryEscape(sender) + "&to=" + url.QueryEscape(receiver) + "&templateId=" + url.QueryEscape(templateId)
	if templateParas != "" {
		param += "&templateParas=" + url.QueryEscape(templateParas)
	}
	if statusCallBack != "" {
		param += "&statusCallback=" + url.QueryEscape(statusCallBack)
	}
	if signature != "" {
		param += "&signature=" + url.QueryEscape(signature)
	}
	return param
}

func post(url string, param []byte, appInfo core.Signer) (string, error) {
	if param == nil || appInfo == (core.Signer{}) {
		return "", nil
	}

	// 代码样例为了简便，设置了不进行证书校验，请在商用环境自行开启证书校验。
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(param))
	if err != nil {
		return "", err
	}

	// 对请求增加内容格式，固定头域
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 对请求进行HMAC算法签名，并将签名结果设置到Authorization头域。
	appInfo.Sign(req)

	slog.Info(fmt.Sprintln(req.Header))
	// 发送短信请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// 获取短信响应
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
