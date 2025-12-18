package miniox

import (
	"fmt"
	"godir/internal/common/svc"
	"time"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

func BuildBaseUrl(useSSL bool, endpoint string) string {
	// 更新数据库中的封面信息
	protocol := "http"
	if useSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, endpoint)
}

type Credentials struct {
	// S3 Access key ID
	AccessKeyID string

	// S3 Secret Access Key
	SecretAccessKey string

	// S3 Session Token
	SessionToken string

	// Expiration of this credentials - null means no expiration associated
	Expiration time.Time
}

func TemporaryCredentials(cfg *svc.Config, paths []string) (*Credentials, error) {
	stsEndpoint := BuildBaseUrl(cfg.MinIO.UseSSL, cfg.MinIO.Endpoint) // STS服务端点，通常与MinIO相同

	opts := credentials.STSAssumeRoleOptions{
		AccessKey: cfg.MinIO.AccessKeyID,                              // 有权限调用STS的Access Key
		SecretKey: cfg.MinIO.SecretAccessKey,                          // 对应的Secret Key
		Location:  cfg.MinIO.Region,                                   // 区域，MinIO中可自定义
		RoleARN:   fmt.Sprintf("arn:aws:s3:::%s/*", cfg.MinIO.Bucket), // test是桶名称
	}
	// for i, v := range paths {
	// 	paths[i] = fmt.Sprintf("arn:aws:s3:::%s/%s", cfg.MinIO.Bucket, strings.TrimPrefix(v, "/"))
	// }
	policy := `{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": [
							"*"
						],
						"Resource": [
							"arn:aws:s3:::*"
						]
					}
				]
			}`
	opts.Policy = policy
	creds, err := credentials.NewSTSAssumeRole(stsEndpoint, opts)
	if err != nil {
		return nil, err
	}

	value, err := creds.GetWithContext(&credentials.CredContext{})
	if err != nil {
		return nil, err
	}
	return &Credentials{
		AccessKeyID:     value.AccessKeyID,
		SecretAccessKey: value.SecretAccessKey,
		SessionToken:    value.SessionToken,
		Expiration:      value.Expiration,
	}, nil
}
