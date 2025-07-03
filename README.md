## 编译可执行文件

编译生成 Linux 下可执行文件

```
## powershell 下通过配置环境变量编译linux可执行文件
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build .\cmd\app.go

# 取消环境变量
$env:GOOS = ""
$env:GOARCH = ""
```

启动选项

| 选项 | 默认值                        | 用途     |
|----|----------------------------|--------|
| -c | `/opt/webhook/config.yaml` | 指定配置文件 |

配置文件默认采用 `/opt/webhook/config.yaml`，可使用 `-c` 指定配置文件，如:

```shell
./app -c config.yaml
```

## 生成Docker镜像

Dockerfile

```Dockerfile
FROM alpine:3.18.3
# 维护者信息
LABEL authors="lglaboy" \
      description="huawei sms webhook,Compatible with Legacy Alerting&New alerts"

# 设置时区，需要安装tzdata
ENV TZ=Asia/Shanghai

# apk add 提速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk add tzdata --no-cache

WORKDIR /opt/webhook

# app 在宿主机上拥有可执行权限，不需要再次授权
COPY app config.yaml /opt/webhook/

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s \
CMD pgrep app || exit 1

ENTRYPOINT ["/opt/webhook/app"]
```

编译镜像

```
docker build -t huawei-sms-webhook .
```

## 启动服务

配置文件示例

```yaml
# 服务启动端口
server:
  port: 8080

# 华为云配置信息
huawei:
  # App Key
  app_key: "apigateway_sdk_demo_key"
  # App Secret
  app_secret: "apigateway_sdk_demo_secret"
  # APP接入地址(在控制台"应用管理"页面获取)+接口访问URI
  api_address: "https://smsapi.cn-south-1.myhuaweicloud.com:443/sms/batchSendSms/v1"
  # 国内短信签名通道号
  sender: "******"
  # 模板ID
  template_id: "******"

  # 条件必填,国内短信关注,当templateId指定的模板类型为通用模板时生效且必填,必须是已审核通过的,与模板类型一致的签名名称
  # 签名名称
  signature: "**信息"
  # 选填,短信状态报告接收地址,推荐使用域名,为空或者不填表示不接收状态报告
  status_call_back: ""

# 告警方式
alerts:
  # 华为云 sms 告警
  - type: huawei-msgsms
    enabled: on
  # 内蒙古医科大学附属医院-短信平台
  - type: nmgykdxfsyy-sms
    enabled: off

# 接收方
receivers:
  - name: 用户1
    # 电话 未添加国家码 +86 ,程序会自动添加 +86
    contact_numbers:
  #      - 12345678901
  #      - +8601234567890
  - name: 用户2
    contact_numbers:
#      - "+8611122220000"
#      - 22233334444
```

操作示例

```bash
# 创建目录
mkdir -p /opt/huawei-sms-webhook

# 获取默认示例配置
docker run -itd --name "huawei-sms-webhook" --entrypoint /bin/bash huawei-sms-webhook
docker cp huawei-sms-webhook:/opt/webhook/config.yaml /opt/huawei-sms-webhook/config.yaml
docker rm -f huawei-sms-webhook

# 编辑配置文件，根据自身调整
vim /opt/huawei-sms-webhook/config.yaml 

# 启动docker容器
docker run -itd \
--name huawei-sms-webhook \
-p 18080:8080 \
--restart=unless-stopped \
--cpus 1 -m 1G \
--log-opt max-size=512m \
--log-opt max-file=3 \
-v /opt/huawei-sms-webhook/config.yaml:/opt/webhook/config.yaml \
huawei-sms-webhook
```

## Grafana配置webhook告警方式

grafana -> Alerting -> Contact points -> 添加一个联络点

联络点配置：

- Name: 自定义
- Integration: Webhook
- URL: `http://192.168.*.*:18080/webhook`



