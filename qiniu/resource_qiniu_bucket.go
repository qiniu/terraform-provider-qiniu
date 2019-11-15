package qiniu

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	qiniu_client "github.com/qiniu/api.v7/v7/client"
	qiniu_storage "github.com/qiniu/api.v7/v7/storage"
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
				Optional:    true,
				Description: "Privately access to the bucket",
			},
			"image_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Image Source URL",
				ValidateFunc: validateURL,
			},
			"image_host": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Image Source Host",
				ValidateFunc: validateHost,
			},
			"lifecycle_rules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Bucket Lifecycle Rules",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "Rule name",
							ForceNew:     true,
							ValidateFunc: validateLifecycleRuleName,
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Rule for object name prefix",
						},
						"to_line_after_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "To line after days",
						},
						"delete_after_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Delete after days",
						},
					},
				},
			},
			"anti_leech_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Anti leech mode",
				ValidateFunc: validateAntiLeechMode,
			},
			"allow_empty_referer": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Allow empty referer",
			},
			"referer_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Referer pattern",
			},
			"only_enable_anti_leech_for_cdn": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Only enable anti leech mode for CDN",
			},
		},
		Create: resourceCreateQiniuBucket,
		Read:   resourceReadQiniuBucket,
		Update: resourceUpdateQiniuBucket,
		Delete: resourceDeleteQiniuBucket,
		Exists: resourceExistsQiniuBucket,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCreateQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	var (
		referAntiLeechConfig qiniu_storage.ReferAntiLeechConfig
		v                    interface{}
		ok                   bool
	)
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Get("name").(string)
	regionID := qiniu_storage.RegionID(d.Get("region_id").(string))
	if err = bucketManager.CreateBucket(bucketName, regionID); err != nil {
		return
	}
	if v, ok = d.GetOk("private"); ok {
		if v.(bool) {
			if err = bucketManager.MakeBucketPrivate(bucketName); err != nil {
				return
			}
		}
	}
	if v, ok = d.GetOk("image_url"); ok {
		imageURL := v.(string)
		if v, ok = d.GetOk("image_host"); ok {
			imageHost := v.(string)
			if err = bucketManager.SetImageWithHost(imageURL, bucketName, imageHost); err != nil {
				return
			}
		} else {
			if err = bucketManager.SetImage(imageURL, bucketName); err != nil {
				return
			}
		}
	}
	if v, ok = d.GetOk("lifecycle_rules"); ok {
		set := v.(*schema.Set)
		for _, r := range set.List() {
			var lifeCycleRule qiniu_storage.BucketLifeCycleRule

			rule := r.(map[string]interface{})
			if v, ok = rule["name"]; ok {
				lifeCycleRule.Name = v.(string)
			}
			if v, ok = rule["prefix"]; ok {
				lifeCycleRule.Prefix = v.(string)
			}
			if v, ok = rule["delete_after_days"]; ok {
				lifeCycleRule.DeleteAfterDays = v.(int)
			}
			if v, ok = rule["to_line_after_days"]; ok {
				lifeCycleRule.ToLineAfterDays = v.(int)
			}
			if err = bucketManager.AddBucketLifeCycleRule(bucketName, &lifeCycleRule); err != nil {
				return
			}
		}
	}
	if v, ok = d.GetOk("anti_leech_mode"); ok {
		switch v.(string) {
		case "":
			referAntiLeechConfig.Mode = 0
		case "whitelist":
			referAntiLeechConfig.Mode = 1
		case "blacklist":
			referAntiLeechConfig.Mode = 2
		default:
			err = errors.New("\"anti_leech_mode\" contains invalid mode")
			return
		}
	}
	if referAntiLeechConfig.Mode > 0 {
		if v, ok = d.GetOk("allow_empty_referer"); ok {
			referAntiLeechConfig.AllowEmptyReferer = v.(bool)
		}
		if v, ok = d.GetOk("referer_pattern"); ok {
			referAntiLeechConfig.Pattern = v.(string)
		}
		if v, ok = d.GetOk("only_enable_anti_leech_for_cdn"); ok {
			referAntiLeechConfig.EnableSource = !v.(bool)
		}
		if err = bucketManager.SetReferAntiLeechMode(bucketName, &referAntiLeechConfig); err != nil {
			return
		}
	}
	d.SetId(bucketName)
	return resourceReadQiniuBucket(d, m)
}

func resourceReadQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	var (
		bucketInfo     qiniu_storage.BucketInfo
		lifeCycleRules []qiniu_storage.BucketLifeCycleRule
	)

	bucketManager := m.(*Client).BucketManager
	bucketName := d.Id()
	bucketInfo, err = bucketManager.GetBucketInfo(bucketName)
	if err == nil {
		lifeCycleRules, err = bucketManager.GetBucketLifeCycleRule(bucketName)
	}

	if err != nil {
		if IsResourceNotFound(err) {
			d.SetId("")
			return nil
		} else {
			return err
		}
	}
	d.Set("name", bucketName)
	d.Set("region_id", bucketInfo.Region)
	d.Set("private", bucketInfo.IsPrivate())
	d.Set("image_url", bucketInfo.Source)
	d.Set("image_host", bucketInfo.Host)
	d.Set("lifecycle_rules", lifeCycleRules)

	switch bucketInfo.AntiLeechMode {
	case 0:
		d.Set("anti_leech_mode", "")
		d.Set("referer_pattern", "")
		d.Set("allow_empty_referer", "")
		d.Set("only_enable_anti_leech_for_cdn", "")
	case 1:
		d.Set("anti_leech_mode", "whitelist")
		d.Set("referer_pattern", strings.Join(bucketInfo.ReferWl, ";"))
		d.Set("allow_empty_referer", bucketInfo.NoRefer)
		d.Set("only_enable_anti_leech_for_cdn", !bucketInfo.EnableSource)
	case 2:
		d.Set("anti_leech_mode", "blacklist")
		d.Set("referer_pattern", strings.Join(bucketInfo.ReferBl, ";"))
		d.Set("allow_empty_referer", bucketInfo.NoRefer)
		d.Set("only_enable_anti_leech_for_cdn", !bucketInfo.EnableSource)
	default:
		err = errors.New("\"anti_leech_mode\" returns server contains invalid mode")
		return
	}
	return nil
}

func resourceUpdateQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	if err = resourcePartialUpdateQiniuBucket(d, m); err != nil {
		return
	}
	return resourceReadQiniuBucket(d, m)
}

func resourcePartialUpdateQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	var (
		bucketManager = m.(*Client).BucketManager
		bucketName    = d.Id()
		v             interface{}
		ok            bool
	)

	if d.HasChange("private") {
		if d.Get("private").(bool) {
			if err = bucketManager.MakeBucketPrivate(bucketName); err != nil {
				return
			}
		} else {
			if err = bucketManager.MakeBucketPublic(bucketName); err != nil {
				return
			}
		}
	}

	if d.HasChange("image_url") || d.HasChange("image_host") {
		if err = bucketManager.UnsetImage(bucketName); err != nil {
			return
		}
		if v, ok = d.GetOk("image_url"); ok {
			imageURL := v.(string)
			if v, ok = d.GetOk("image_host"); ok {
				imageHost := v.(string)
				if err = bucketManager.SetImageWithHost(imageURL, bucketName, imageHost); err != nil {
					return
				}
			} else {
				if err = bucketManager.SetImage(imageURL, bucketName); err != nil {
					return
				}
			}
		}
	}
	if d.HasChange("lifecycle_rules") {
		var (
			ruleName                 string
			newRule                  qiniu_storage.BucketLifeCycleRule
			oldRulesList             []qiniu_storage.BucketLifeCycleRule
			oldRulesMap, newRulesMap map[string]qiniu_storage.BucketLifeCycleRule
		)
		if oldRulesList, err = bucketManager.GetBucketLifeCycleRule(bucketName); err != nil {
			return
		} else {
			oldRulesMap = make(map[string]qiniu_storage.BucketLifeCycleRule, len(oldRulesList))
			for i := range oldRulesList {
				oldRulesMap[oldRulesList[i].Name] = oldRulesList[i]
			}
		}
		if v, ok = d.GetOk("lifecycle_rules"); ok {
			set := v.(*schema.Set)
			newRulesMap = make(map[string]qiniu_storage.BucketLifeCycleRule, set.Len())

			for _, r := range set.List() {
				var newRule qiniu_storage.BucketLifeCycleRule

				rule := r.(map[string]interface{})
				if v, ok = rule["name"]; ok {
					newRule.Name = v.(string)
				}
				if v, ok = rule["prefix"]; ok {
					newRule.Prefix = v.(string)
				}
				if v, ok = rule["delete_after_days"]; ok {
					newRule.DeleteAfterDays = v.(int)
				}
				if v, ok = rule["to_line_after_days"]; ok {
					newRule.ToLineAfterDays = v.(int)
				}
				newRulesMap[newRule.Name] = newRule
			}
		}
		for ruleName, _ = range oldRulesMap {
			if newRule, ok = newRulesMap[ruleName]; ok {
				if err = bucketManager.UpdateBucketLifeCycleRule(bucketName, &newRule); err != nil {
					return err
				}
			} else {
				if err = bucketManager.DelBucketLifeCycleRule(bucketName, ruleName); err != nil {
					return err
				}
			}
		}
		for ruleName, newRule = range newRulesMap {
			if _, ok = oldRulesMap[ruleName]; ok {
				if err = bucketManager.AddBucketLifeCycleRule(bucketName, &newRule); err != nil {
					return err
				}
			}
		}
	}

	if d.HasChange("anti_leech_mode") || d.HasChange("referer_pattern") || d.HasChange("allow_empty_referer") || d.HasChange("only_enable_anti_leech_for_cdn") {
		var referAntiLeechConfig qiniu_storage.ReferAntiLeechConfig
		if v, ok = d.GetOk("anti_leech_mode"); ok {
			switch v.(string) {
			case "":
				referAntiLeechConfig.Mode = 0
			case "whitelist":
				referAntiLeechConfig.Mode = 1
			case "blacklist":
				referAntiLeechConfig.Mode = 2
			default:
				err = errors.New("\"anti_leech_mode\" contains invalid mode")
				return
			}
		}
		if v, ok = d.GetOk("allow_empty_referer"); ok {
			referAntiLeechConfig.AllowEmptyReferer = v.(bool)
		}
		if v, ok = d.GetOk("referer_pattern"); ok {
			referAntiLeechConfig.Pattern = v.(string)
		}
		if v, ok = d.GetOk("only_enable_anti_leech_for_cdn"); ok {
			referAntiLeechConfig.EnableSource = !v.(bool)
		}
		if err = bucketManager.SetReferAntiLeechMode(bucketName, &referAntiLeechConfig); err != nil {
			return
		}
	}

	return nil
}

func resourceDeleteQiniuBucket(d *schema.ResourceData, m interface{}) (err error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Id()

	if err = bucketManager.DropBucket(bucketName); err != nil {
		if !IsResourceNotFound(err) {
			return err
		}
	}
	d.SetId("")
	return nil
}

func resourceExistsQiniuBucket(d *schema.ResourceData, m interface{}) (bool, error) {
	bucketManager := m.(*Client).BucketManager
	bucketName := d.Id()
	if _, err := bucketManager.GetBucketInfo(bucketName); err == nil {
		return true, nil
	} else if IsResourceNotFound(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsResourceNotFound(err error) bool {
	if qiniuErr, ok := err.(*qiniu_client.ErrorInfo); ok {
		return qiniuErr.HttpCode() == HTTP_STATUS_RESOURCE_NOT_FOUND || qiniuErr.HttpCode() == HTTP_STATUS_BUCKET_NOT_FOUND
	}
	return false
}
