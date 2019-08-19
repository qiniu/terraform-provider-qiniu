.PHONY: build

build: bin/terraform-provider-qiniu

bin/terraform-provider-qiniu:
	go build -o bin/terraform-provider-qiniu
