package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/qiniu/terraform-provider-qiniu/qiniu"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qiniu.Provider,
	})
}
