.PHONY: build test bin/terraform-provider-qiniu clean

build: bin/terraform-provider-qiniu
test:
	GO111MODULE=on go test -v ./qiniu/... -args -ginkgo.failFast -ginkgo.progress -ginkgo.v -ginkgo.trace -test.parallel 1

bin/terraform-provider-qiniu:
	GO111MODULE=on go build -o bin/terraform-provider-qiniu

clean:
	rm -rf bin/terraform-provider-qiniu
