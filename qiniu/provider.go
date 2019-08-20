package qiniu

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qiniu_auth "github.com/qiniu/api.v7/auth"
	qiniu_storage "github.com/qiniu/api.v7/storage"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("QINIU_ACCESS_KEY", ""),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("QINIU_SECRET_KEY", ""),
			},
			"use_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("QINIU_USE_HTTPS", false),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"qiniu_bucket": resourceQiniuBucket(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"qiniu_buckets": dataSourceQiniuBuckets(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		storageConfig qiniu_storage.Config
		auth          = qiniu_auth.New(d.Get("access_key").(string), d.Get("secret_key").(string))
	)

	storageConfig.UseHTTPS = d.Get("use_https").(bool)

	return &Client{
		Auth:          auth,
		BucketManager: qiniu_storage.NewBucketManager(auth, &storageConfig),
	}, nil
}
