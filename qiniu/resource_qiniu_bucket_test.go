package qiniu_test

import (
	"fmt"

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
