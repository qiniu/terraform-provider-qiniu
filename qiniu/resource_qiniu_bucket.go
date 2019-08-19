package qiniu

import (
	"github.com/hashicorp/terraform/helper/schema"
	qiniu_client "github.com/qiniu/api.v7/client"
	qiniu_storage "github.com/qiniu/api.v7/storage"
)

const (
	HTTP_STATUS_RESOURCE_NOT_FOUND = 612
	HTTP_STATUS_BUCKET_NOT_FOUND   = 631
)

func resourceQiniuBucket() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the bucket",
				ForceNew:     true,
				ValidateFunc: validateBucketName,
			},
			"region_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The region id of the bucket",
				ForceNew:     true,
				ValidateFunc: validateRegionID,
			},
			"private": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Privately access to the bucket",
			},
		},
		Create: resourceCreateQiniuBucket,
		Read:   resourceReadQiniuBucket,
		Delete: resourceDeleteQiniuBucket,
		Exists: resourceExistsQiniuBucket,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Get("name").(string)
	regionID := qiniu_storage.RegionID(d.Get("region_id").(string))
	if err = bucketManager.CreateBucket(bucketName, regionID); err != nil {
		return
	}
	if d.Get("private").(bool) {
		if err = bucketManager.MakeBucketPrivate(bucketName); err != nil {
			return
		}
	}
	d.SetId(bucketName)
	return nil
}

func resourceReadQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	var bucketInfo qiniu_storage.BucketInfo

	bucketManager := m.(*Client).BucketManager
	bucketName := d.Id()
	bucketInfo, err = bucketManager.GetBucketInfo(bucketName)

	if err != nil {
		if IsBucketNotFound(err) {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	d.SetId(bucketName)
	d.Set("name", bucketName)
	d.Set("region_id", bucketInfo.Region)
	d.Set("private", bucketInfo.IsPrivate())
	return nil
}

func resourceDeleteQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Id()

	if err = bucketManager.DropBucket(bucketName); err != nil {
		if !IsBucketNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}

func resourceExistsQiniuBucket(d *schema.ResourceData, m interface{}) (bool, error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Get("name").(string)
	if _, err := bucketManager.GetBucketInfo(bucketName); err == nil {
		return true, nil
	} else if IsBucketNotFound(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsBucketNotFound(err error) bool {
	if qiniuErr, ok := err.(*qiniu_client.ErrorInfo); ok {
		return qiniuErr.Code == HTTP_STATUS_RESOURCE_NOT_FOUND || qiniuErr.Code == HTTP_STATUS_BUCKET_NOT_FOUND
	}
	return false
}
