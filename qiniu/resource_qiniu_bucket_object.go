package qiniu

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	qiniu_storage "github.com/qiniu/api.v7/v7/storage"
)

func resourceQiniuBucketObject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The bucket to contain the object",
				ForceNew:     true,
				ValidateFunc: validateBucketName,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The object key",
				ForceNew:    true,
			},
			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content"},
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source"},
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
		Create: resourcePutQiniuBucketObject,
		Read:   resourceReadQiniuBucketObject,
		Update: resourcePutQiniuBucketObject,
		Delete: resourceDeleteQiniuBucketObject,
		Exists: resourceExistsQiniuBucketObject,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourcePutQiniuBucketObject(d *schema.ResourceData, m interface{}) (err error) {
	var (
		reader io.ReaderAt
		size   int64
		bucket = d.Get("bucket").(string)
		key    = d.Get("key").(string)
		entry  = getEntryFromBucketNameAndKey(bucket, key)
	)

	putPolicy := qiniu_storage.PutPolicy{Scope: entry, DetectMime: 1}
	uploadToken := putPolicy.UploadToken(m.(*Client).Auth)

	if source, ok := d.GetOk("source"); ok {
		var (
			file     *os.File
			fileInfo os.FileInfo
		)
		if file, err = os.Open(source.(string)); err != nil {
			return
		}
		if fileInfo, err = file.Stat(); err != nil {
			return
		}
		reader = file
		size = fileInfo.Size()
	} else if content, ok := d.GetOk("content"); ok {
		r := strings.NewReader(content.(string))
		reader = r
		size = r.Size()
	} else {
		return errors.New("Neither \"source\" nor \"content\" is specified")
	}
	if err = m.(*Client).ResumeUploader.Put(
		context.Background(), nil, uploadToken, key, reader, size, &qiniu_storage.RputExtra{TryTimes: 5},
	); err != nil {
		return
	}
	d.SetId(entry)
	return resourceReadQiniuBucketObject(d, m)
}

func resourceReadQiniuBucketObject(d *schema.ResourceData, m interface{}) error {
	var bucketManager = m.(*Client).BucketManager

	if fileInfo, err := bucketManager.Stat(getBucketNameAndKeyFromEntry(d.Id())); err != nil {
		return err
	} else {
		d.Set("content_type", fileInfo.MimeType)
		d.Set("content_length", fileInfo.Fsize)
		d.Set("content_etag", fileInfo.Hash)
		return nil
	}
}

func resourceDeleteQiniuBucketObject(d *schema.ResourceData, m interface{}) error {
	var bucketManager = m.(*Client).BucketManager

	if err := bucketManager.Delete(getBucketNameAndKeyFromEntry(d.Id())); err != nil {
		if IsResourceNotFound(err) {
			return nil
		}
		return err
	}
	d.SetId("")
	return nil
}

func resourceExistsQiniuBucketObject(d *schema.ResourceData, m interface{}) (bool, error) {
	var bucketManager = m.(*Client).BucketManager

	if _, err := bucketManager.Stat(getBucketNameAndKeyFromEntry(d.Id())); err != nil {
		if IsResourceNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getBucketNameAndKeyFromEntry(entry string) (bucket string, key string) {
	parts := strings.SplitN(entry, ":", 2)
	return parts[0], parts[1]
}

func getEntryFromBucketNameAndKey(bucket string, key string) (entry string) {
	return bucket + ":" + key
}
