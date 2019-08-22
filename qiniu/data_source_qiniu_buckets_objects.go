package qiniu

import (
	"github.com/hashicorp/terraform/helper/schema"
	qiniu_storage "github.com/qiniu/api.v7/storage"
)

func dataSourceQiniuBucketsObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQiniuBucketsObjectsRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateBucketName,
				ForceNew:     true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePositiveInt,
				ForceNew:     true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"key_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"content_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"content_etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceQiniuBucketsObjectsRead(d *schema.ResourceData, m interface{}) (err error) {
	var (
		bucketManager       = m.(*Client).BucketManager
		bucket              = d.Get("bucket").(string)
		prefix              string
		marker              string
		entries, allEntries []qiniu_storage.ListItem
		totalLimit          int = (1 << 31) - 1
		hasNext                 = true
	)

	if v, ok := d.GetOk("prefix"); ok {
		prefix = v.(string)
	}
	if v, ok := d.GetOk("limit"); ok {
		totalLimit = v.(int)
	}

	for hasNext && totalLimit > 0 {
		limit := 1000
		if totalLimit < 1000 {
			limit = totalLimit
		}
		if entries, _, marker, hasNext, err = bucketManager.ListFiles(bucket, prefix, "", marker, limit); err != nil {
			return err
		}
		allEntries = append(allEntries, entries...)
		totalLimit -= len(entries)
	}
	return dataSourceQiniuBucketsObjectsAttributes(d, bucket, allEntries)
}

func dataSourceQiniuBucketsObjectsAttributes(d *schema.ResourceData, bucket string, entries []qiniu_storage.ListItem) (err error) {
	var (
		ids      = make([]string, 0, len(entries))
		keys     = make([]string, 0, len(entries))
		keyInfos = make([]map[string]interface{}, 0, len(entries))
	)

	for _, entry := range entries {
		attributes := map[string]interface{}{
			"key":            entry.Key,
			"content_etag":   entry.Hash,
			"content_length": entry.Fsize,
			"content_type":   entry.MimeType,
		}
		ids = append(ids, getEntryFromBucketNameAndKey(bucket, entry.Key))
		keys = append(keys, entry.Key)
		keyInfos = append(keyInfos, attributes)
	}
	d.SetId(dataResourceIdHash(ids))
	if err = d.Set("keys", keys); err != nil {
		return
	}
	if err = d.Set("key_infos", keyInfos); err != nil {
		return
	}
	return
}
