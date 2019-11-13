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
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
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
}
                `,
				ExpectError: regexp.MustCompile("must not be empty"),
			}},
		})
	})

	It("should reject too short qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "ab"
    region_id = "z2"
}
                `,
				ExpectError: regexp.MustCompile("must not be shorter than 3 characters"),
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
    name = "valid-name"
    region_id = "z100"
}
                `,
				ExpectError: regexp.MustCompile("is invalid"),
			}},
		})
	})

	It("should reject invalid qiniu image_url", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name"
    region_id = "z1"
    image_url = "www.qiniu.com"
}
                `,
				ExpectError: regexp.MustCompile("must be valid url"),
			}},
		})
	})

	It("should reject invalid qiniu image_host", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name"
    region_id = "z1"
    image_url = "http://www.qiniu.com"
    image_host = "http://www.qiniu.com"
}
                `,
				ExpectError: regexp.MustCompile("must be valid host"),
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
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
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
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
    image_url = "http://www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", "http://www.qiniu.com"),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
    image_url = "http://portal.qiniu.io"
    image_host = "www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", "http://portal.qiniu.io"),
					resource.TestCheckResourceAttr(resourceID, "image_host", "www.qiniu.com"),
				),
			}, {
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
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
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
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "image_url", ""),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "image_host", ""),
					testCheckQiniuBucketItemExists("qiniu_bucket.private-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "name", "bucket-test-2-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "private", "true"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "image_url", ""),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "imaged-bucket" {
    name = "bucket-test-3-terraform"
    region_id = "z2"
    image_url = "http://qiniu.io"
}

resource "qiniu_bucket" "imaged-bucket-with-host" {
    name = "bucket-test-4-terraform"
    region_id = "z2"
    image_url = "http://qiniu.io"
    image_host = "www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.imaged-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "name", "bucket-test-3-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "image_url", "http://qiniu.io"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "image_host", ""),
					testCheckQiniuBucketItemExists("qiniu_bucket.imaged-bucket-with-host"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "name", "bucket-test-4-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "image_url", "http://qiniu.io"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "image_host", "www.qiniu.com"),
				),
			}},
		})
	})
})
