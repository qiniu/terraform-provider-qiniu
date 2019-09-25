package qiniu_test

import (
	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("dataSourceQiniuRegions", func() {
	It("should list qiniu regions", func() {
		resource.Test(MakeT("TestListQiniuRegions"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
data "qiniu_regions" "all" {
}

data "qiniu_regions" "china" {
    description_regex = "China"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.qiniu_regions.all", "regions.#", "5"),
					resource.TestCheckResourceAttr("data.qiniu_regions.all", "region_ids.#", "5"),
					resource.TestCheckResourceAttr("data.qiniu_regions.china", "regions.#", "3"),
					resource.TestCheckResourceAttr("data.qiniu_regions.china", "region_ids.#", "3"),
				),
			}},
		})
	})
})
