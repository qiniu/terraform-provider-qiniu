module github.com/qiniu/terraform-provider-qiniu

go 1.12

require (
	github.com/hashicorp/terraform v0.12.9
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/qiniu/api.v7 v0.0.0-20190520053455-bea02cd22bf4
)

replace github.com/qiniu/api.v7 => ./qiniu/sdk
