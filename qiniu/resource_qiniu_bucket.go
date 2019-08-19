package qiniu

import (
	"github.com/hashicorp/terraform/helper/schema"
	qiniu_storage "github.com/qiniu/api.v7/storage"
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
				Description: "Privately access to the bucket",
			},
		},
		Create: resourceCreateQiniuBucket,
		// Read:   resourceReadQiniuBucket,
		// Update: resourceUpdateQiniuBucket,
		// Delete: resourceDeleteQiniuBucket,
		Exists: resourceExistsQiniuBucket,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
	}
}

func resourceCreateQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Get("name").(string)
	regionID := qiniu_storage.RegionID(d.Get("region_id").(string))
	if err = bucketManager.CreateBucket(bucketName, regionID); err != nil {
		return
	}
	d.SetId(bucketName)
	return nil
}

func resourceExistsQiniuBucket(d *schema.ResourceData, m interface{}) (bool, error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Get("name").(string)
	_, err := bucketManager.GetBucketInfo(bucketName)
	return err == nil, nil
}
