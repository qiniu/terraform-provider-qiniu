.PHONY: build test bin/terraform-provider-qiniu clean

build: bin/terraform-provider-qiniu
test:
	GO111MODULE=on go test -v ./qiniu/...

bin/terraform-provider-qiniu:
	GO111MODULE=on go build -o bin/terraform-provider-qiniu

clean:
	rm -rf bin/terraform-provider-qiniu
