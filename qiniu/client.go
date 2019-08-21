package qiniu

import (
	qiniu_auth "github.com/qiniu/api.v7/auth"
	qiniu_storage "github.com/qiniu/api.v7/storage"
)

type Client struct {
	Auth           *qiniu_auth.Credentials
	BucketManager  *qiniu_storage.BucketManager
	ResumeUploader *qiniu_storage.ResumeUploader
}
