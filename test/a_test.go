package test

import (
	"fmt"
	"testing"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestRunAppExe(t *testing.T) {
	// 1. 创建STS凭证提供者
	// endpoint := "your-minio-server:9000" // MinIO服务地址
	stsEndpoint := "http://127.0.0.1:9000" // STS服务端点，通常与MinIO相同
	creds, err := credentials.NewSTSAssumeRole(stsEndpoint, credentials.STSAssumeRoleOptions{
		AccessKey: "sts-user",            // 有权限调用STS的Access Key
		SecretKey: "11111111",            // 对应的Secret Key
		Location:  "mengfanbing",         // 区域，MinIO中可自定义
		RoleARN:   "arn:aws:s3:::test/*", // test是桶名称
		Policy: `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": [
							"s3:GetBucketLocation",
							"s3:GetObject",
							"s3:PutObject",
							"s3:DeleteObject"
						],
						"Resource": [
							"arn:aws:s3:::test/*"
						]
					}
				]
			}`,
		// DurationSeconds: time.,
	})
	if err != nil {
		fmt.Println("创建STS凭证提供者失败:", err)
		return
	}

	value, err := creds.GetWithContext(&credentials.CredContext{})
	// value, err := creds.Get()
	if err != nil {
		fmt.Println("GET():", err)
		return
	}
	fmt.Println("token:", value.SessionToken)
}
