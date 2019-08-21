package qiniu_test

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("resourceQiniuBucket", func() {
	It("should create qiniu bucket", func() {
		resourceID := "qiniu_bucket.basic_bucket"
		resource.Test(MakeT("TestCreateQiniuBucket"), resource.TestCase{
			PreCheck:      testPreCheck,
			IDRefreshName: resourceID,
			Providers:     providers,
			CheckDestroy:  testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "basic_bucket" {
    name = "basic-test-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "basic-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "true"),
				),
			}},
		})
	})

	It("should reject invalid qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "invalid(bucket)name"
    region_id = "z2"
    private = true
}
                `,
				ExpectError: regexp.MustCompile("must not contain invalid characters"),
			}},
		})
	})

	It("should reject empty qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = ""
    region_id = "z2"
    private = true
}
                `,
				ExpectError: regexp.MustCompile("must not be empty"),
			}},
		})
	})

	It("should reject too long qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglong"
    region_id = "z2"
    private = true
}
                `,
				ExpectError: regexp.MustCompile("must not be longer than 63 characters"),
			}},
		})
	})

	It("should reject invalid qiniu region id", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid_name"
    region_id = "z100"
    private = true
}
                `,
				ExpectError: regexp.MustCompile("is invalid"),
			}},
		})
	})

	It("should update qiniu bucket", func() {
		resourceID := "qiniu_bucket.update_bucket"
		resource.Test(MakeT("TestUpdateQiniuBucket"), resource.TestCase{
			PreCheck:      testPreCheck,
			IDRefreshName: resourceID,
			Providers:     providers,
			CheckDestroy:  testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "true"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
				),
			}},
		})
	})

	It("should create qiniu buckets", func() {
		resource.Test(MakeT("TestCreateQiniuBuckets"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "public-bucket" {
    name = "bucket-test-1-terraform"
    region_id = "z2"
    private = false
}

resource "qiniu_bucket" "private-bucket" {
    name = "bucket-test-2-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.public-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "name", "bucket-test-1-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "private", "false"),
					testCheckQiniuBucketItemExists("qiniu_bucket.private-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "name", "bucket-test-2-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "private", "true"),
				),
			}},
		})
	})
})
