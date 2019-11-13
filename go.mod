module github.com/qiniu/terraform-provider-qiniu

go 1.12

require (
	github.com/hashicorp/terraform v0.12.9
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/api.v7/v7 v7.4.0
	github.com/qiniu/x v7.0.8+incompatible // indirect
	qiniupkg.com/x v7.0.8+incompatible // indirect
)

replace github.com/qiniu/api.v7/v7 => ./qiniu/sdk
