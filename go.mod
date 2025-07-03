module huawei-msgsms-webhook

require (
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	gopkg.in/yaml.v2 v2.4.0
	huaweicloud.com/apig/signer v0.0.0
)

replace huaweicloud.com/apig/signer v0.0.0 => ./pkg/APIGW-go-sdk-2.0.2/core

go 1.20
