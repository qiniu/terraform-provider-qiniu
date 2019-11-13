package qiniu

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	qiniu_auth "github.com/qiniu/api.v7/v7/auth"
	qiniu_client "github.com/qiniu/api.v7/v7/client"
	qiniu_storage "github.com/qiniu/api.v7/v7/storage"
)

func init() {
	qiniu_client.SetAppName("terraform-provider-qiniu")
}

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
			"central_rs_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_CENTRAL_RS_URL", ""),
				ValidateFunc: validateURL,
			},
			"rs_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_RS_URL", ""),
				ValidateFunc: validateURL,
			},
			"rsf_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_RSF_URL", ""),
				ValidateFunc: validateURL,
			},
			"up_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_UP_URL", ""),
				ValidateFunc: validateURL,
			},
			"api_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_API_URL", ""),
				ValidateFunc: validateURL,
			},
			"io_url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("QINIU_IO_URL", ""),
				ValidateFunc: validateURL,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"qiniu_bucket":        resourceQiniuBucket(),
			"qiniu_bucket_object": resourceQiniuBucketObject(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"qiniu_buckets":         dataSourceQiniuBuckets(),
			"qiniu_regions":         dataSourceQiniuRegions(),
			"qiniu_buckets_objects": dataSourceQiniuBucketsObjects(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var (
		storageConfig qiniu_storage.Config
		auth          = qiniu_auth.New(d.Get("access_key").(string), d.Get("secret_key").(string))
		v             interface{}
		ok            bool
	)

	storageConfig.UseHTTPS = d.Get("use_https").(bool)

	if v, ok = d.GetOk("central_rs_url"); ok {
		storageConfig.CentralRsHost = v.(string)
	}
	if v, ok = d.GetOk("rs_url"); ok {
		storageConfig.RsHost = v.(string)
	}
	if v, ok = d.GetOk("rsf_url"); ok {
		storageConfig.RsfHost = v.(string)
	}
	if v, ok = d.GetOk("api_url"); ok {
		storageConfig.ApiHost = v.(string)
	}
	if v, ok = d.GetOk("io_url"); ok {
		storageConfig.IoHost = v.(string)
	}
	if v, ok = d.GetOk("up_url"); ok {
		storageConfig.UpHost = v.(string)
	}

	return &Client{
		Auth:           auth,
		BucketManager:  qiniu_storage.NewBucketManager(auth, &storageConfig),
		ResumeUploader: qiniu_storage.NewResumeUploader(&storageConfig),
	}, nil
}
