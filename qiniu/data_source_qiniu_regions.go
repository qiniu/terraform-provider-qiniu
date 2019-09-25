package qiniu

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	qiniu_storage "github.com/qiniu/api.v7/storage"
)

func dataSourceQiniuRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceQiniuRegionsRead,
		Schema: map[string]*schema.Schema{
			"id_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"region_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceQiniuRegionsRead(d *schema.ResourceData, m interface{}) (err error) {
	var regionsInfo []qiniu_storage.RegionInfo

	auth := m.(*Client).Auth
	if regionsInfo, err = qiniu_storage.GetRegionsInfo(auth); err != nil {
		return
	}

	if v, ok := d.GetOk("id_regex"); ok && v.(string) != "" {
		idRegexp := regexp.MustCompile(v.(string))
		allRegionsInfo := regionsInfo
		regionsInfo = make([]qiniu_storage.RegionInfo, 0, len(allRegionsInfo))
		for _, regionInfo := range allRegionsInfo {
			if idRegexp.MatchString(regionInfo.ID) {
				regionsInfo = append(regionsInfo, regionInfo)
			}
		}
	}
	if v, ok := d.GetOk("description_regex"); ok && v.(string) != "" {
		descriptionRegexp := regexp.MustCompile(v.(string))
		allRegionsInfo := regionsInfo
		regionsInfo = make([]qiniu_storage.RegionInfo, 0, len(allRegionsInfo))
		for _, regionInfo := range allRegionsInfo {
			if descriptionRegexp.MatchString(regionInfo.Description) {
				regionsInfo = append(regionsInfo, regionInfo)
			}
		}
	}
	return dataSourceQiniuRegionsAttributes(d, regionsInfo)
}

func dataSourceQiniuRegionsAttributes(d *schema.ResourceData, regionsInfo []qiniu_storage.RegionInfo) (err error) {
	var (
		regionIds   = make([]string, 0, len(regionsInfo))
		regionInfos = make([]map[string]interface{}, 0, len(regionsInfo))
	)
	for _, regionInfo := range regionsInfo {
		attributes := map[string]interface{}{
			"id":          regionInfo.ID,
			"description": regionInfo.Description,
		}

		regionIds = append(regionIds, regionInfo.ID)
		regionInfos = append(regionInfos, attributes)
	}
	d.SetId(dataResourceIdHash(regionIds))
	if err = d.Set("regions", regionInfos); err != nil {
		return
	}
	if err = d.Set("region_ids", regionIds); err != nil {
		return
	}
	return
}
