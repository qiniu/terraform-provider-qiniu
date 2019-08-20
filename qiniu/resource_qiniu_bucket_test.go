package qiniu_test

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	. "github.com/onsi/ginkgo"
	qiniu "github.com/qiniu/terraform-provider-qiniu/qiniu"
)

var _ = Describe("resourceQiniuBucket", func() {
	It("should create qiniu bucket", func() {
		resourceID := "qiniu_bucket.basic_bucket"
		resource.Test(MakeT("TestCreateQiniuBucket"), resource.TestCase{
			PreCheck:      testPreCheck,
			IDRefreshName: resourceID,
			Providers:     providers,
			CheckDestroy:  testCheckQiniuBucketItemDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "basic_bucket" {
    name = "basic-test-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeTestCheckFunc(
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
				ExpectError: regexp.MustCompile("invalid arguments"),
			}},
		})

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
			CheckDestroy:  testCheckQiniuBucketItemDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeTestCheckFunc(
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
				Check: resource.ComposeTestCheckFunc(
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
			CheckDestroy: testCheckQiniuBucketItemDestroy,
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
				Check: resource.ComposeTestCheckFunc(
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

func testCheckQiniuBucketItemExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		bucketName := rs.Primary.ID
		client := qiniuProvider.Meta().(*qiniu.Client)
		_, err := client.BucketManager.GetBucketInfo(bucketName)
		return err
	}
}

func testCheckQiniuBucketItemDestroy(s *terraform.State) (err error) {
	client := qiniuProvider.Meta().(*qiniu.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "qiniu_bucket" {
			continue
		}

		bucketName := rs.Primary.ID
		if _, err = client.BucketManager.GetBucketInfo(bucketName); err == nil {
			return fmt.Errorf("Alert still exists")
		} else if !qiniu.IsBucketNotFound(err) {
			return
		}
	}

	return nil
}
