package qiniu

import (
	"errors"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	qiniu_storage "github.com/qiniu/api.v7/v7/storage"
)

func dataSourceQiniuBuckets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQiniuBucketsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"region_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegionID,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"image_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_page_on": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_age": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"anti_leech_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allow_empty_referer": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"referer_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"only_enable_anti_leech_for_cdn": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceQiniuBucketsRead(d *schema.ResourceData, m interface{}) (err error) {
	var (
		buckets  []qiniu_storage.BucketSummary
		regionId qiniu_storage.RegionID
	)

	bucketManager := m.(*Client).BucketManager
	if v, ok := d.GetOk("region_id"); ok && v.(string) != "" {
		regionId = qiniu_storage.RegionID(v.(string))
	}
	if buckets, err = bucketManager.BucketInfosInRegion(regionId, false); err != nil {
		return
	}
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		nameRegexp := regexp.MustCompile(v.(string))
		allBuckets := buckets
		buckets = make([]qiniu_storage.BucketSummary, 0, len(allBuckets))
		for _, bucket := range allBuckets {
			if nameRegexp.MatchString(bucket.Name) {
				buckets = append(buckets, bucket)
			}
		}
	}
	return dataSourceQiniuBucketsAttributes(d, buckets)
}

func dataSourceQiniuBucketsAttributes(d *schema.ResourceData, buckets []qiniu_storage.BucketSummary) (err error) {
	var (
		ids         = make([]string, 0, len(buckets))
		bucketNames = make([]string, 0, len(buckets))
		bucketInfos = make([]map[string]interface{}, 0, len(buckets))
	)

	for _, bucket := range buckets {
		attributes := map[string]interface{}{
			"name":          bucket.Name,
			"region_id":     bucket.Info.Region,
			"private":       bucket.Info.IsPrivate(),
			"image_url":     bucket.Info.Source,
			"image_host":    bucket.Info.Host,
			"index_page_on": bucket.Info.IndexPageOn(),
			"max_age":       bucket.Info.MaxAge,
		}

		switch bucket.Info.AntiLeechMode {
		case 0:
			attributes["anti_leech_mode"] = ""
			attributes["referer_pattern"] = nil
			attributes["allow_empty_referer"] = nil
			attributes["only_enable_anti_leech_for_cdn"] = nil
		case 1:
			attributes["anti_leech_mode"] = "whitelist"
			attributes["referer_pattern"] = strings.Join(bucket.Info.ReferWl, ";")
			attributes["allow_empty_referer"] = bucket.Info.NoRefer
			attributes["only_enable_anti_leech_for_cdn"] = !bucket.Info.EnableSource
		case 2:
			attributes["anti_leech_mode"] = "blacklist"
			attributes["referer_pattern"] = strings.Join(bucket.Info.ReferBl, ";")
			attributes["allow_empty_referer"] = bucket.Info.NoRefer
			attributes["only_enable_anti_leech_for_cdn"] = !bucket.Info.EnableSource
		default:
			err = errors.New("\"anti_leech_mode\" returns by server contains invalid mode")
			return
		}

		ids = append(ids, bucket.Name)
		bucketNames = append(bucketNames, bucket.Name)
		bucketInfos = append(bucketInfos, attributes)
	}
	d.SetId(dataResourceIdHash(ids))
	if err = d.Set("buckets", bucketInfos); err != nil {
		return
	}
	if err = d.Set("names", bucketNames); err != nil {
		return
	}
	return
}
