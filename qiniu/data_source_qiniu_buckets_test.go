package qiniu_test

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("dataSourceQiniuBuckets", func() {
	It("should list qiniu buckets", func() {
		resource.Test(MakeT("TestCreateAndListQiniuBuckets"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "basic_bucket_1" {
    name = "basic-test-terraform-1-%d"
    region_id = "z2"
    private = true
}

resource "qiniu_bucket" "basic_bucket_2" {
    name = "basic-test-terraform-2-%d"
    region_id = "z1"
    private = false
}

resource "qiniu_bucket" "basic_bucket_3" {
    name = "basic-test-terraform-3-%d"
    region_id = "as0"
    private = true
}
                `, time.Now().UnixNano(), time.Now().UnixNano(), time.Now().UnixNano()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket_1"),
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket_2"),
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket_3"),
				),
			}, {
				Config: `
data "qiniu_buckets" "all" {
    name_regex = "^basic-test-terraform-"
}

data "qiniu_buckets" "z1" {
    name_regex = "^basic-test-terraform-"
    region_id = "z1"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.qiniu_buckets.all", "buckets.#", "3"),
					resource.TestCheckResourceAttr("data.qiniu_buckets.all", "names.#", "3"),
					resource.TestCheckResourceAttr("data.qiniu_buckets.z1", "buckets.#", "1"),
				),
			}},
		})
	})

	It("should verify qiniu buckets filter syntax", func() {
		resource.Test(MakeT("TestVerifyQiniuBucketsFilter"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
data "qiniu_buckets" "all" {
    name_regex = "oo[xx"
}
                `,
				ExpectError: regexp.MustCompile("error parsing regexp"),
			}, {
				Config: `
data "qiniu_buckets" "all" {
    region_id = "z100"
}
                `,
				ExpectError: regexp.MustCompile("is invalid"),
			}},
		})
	})
})
