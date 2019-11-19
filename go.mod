module github.com/qiniu/terraform-provider-qiniu

go 1.12

require (
	cloud.google.com/go v0.48.0 // indirect
	cloud.google.com/go/storage v1.3.0 // indirect
	github.com/aws/aws-sdk-go v1.25.37 // indirect
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/hashicorp/go-hclog v0.10.0 // indirect
	github.com/hashicorp/go-plugin v1.0.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/terraform v0.12.16
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/joho/godotenv v1.3.0
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/posener/complete v1.2.3 // indirect
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/api.v7/v7 v7.4.0
	github.com/qiniu/x v7.0.8+incompatible // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/ulikunitz/xz v0.5.6 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	go.opencensus.io v0.22.2 // indirect
	golang.org/x/crypto v0.0.0-20191117063200-497ca9f6d64f // indirect
	golang.org/x/net v0.0.0-20191118183410-d06c31c94cae // indirect
	golang.org/x/sys v0.0.0-20191118133127-cf1e2d577169 // indirect
	golang.org/x/tools v0.0.0-20191118222007-07fc4c7f2b98 // indirect
	google.golang.org/api v0.14.0 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20191115221424-83cc0476cb11 // indirect
	google.golang.org/grpc v1.25.1 // indirect
	qiniupkg.com/x v7.0.8+incompatible // indirect
)

replace github.com/qiniu/api.v7/v7 => ./qiniu/sdk
