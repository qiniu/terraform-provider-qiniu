package qiniu_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("dataSourceQiniuBucketsObjects", func() {
	It("should list qiniu buckets objects", func() {
		var (
			tmpFiles = make([]*os.File, 3)
			err      error
		)
		for i := 0; i < 3; i++ {
			tmpFiles[i], err = ioutil.TempFile("", "")
			Expect(err).NotTo(HaveOccurred())
			defer func(path string) {
				os.Remove(path)
			}(tmpFiles[i].Name())
			for j := 0; j <= i; j++ {
				_, err = io.WriteString(tmpFiles[i], "hello world")
				Expect(err).NotTo(HaveOccurred())
			}
			Expect(tmpFiles[i].Close()).To(Succeed())
		}

		resource.Test(MakeT("TestCreateAndListQiniuBucketsObjects"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket_object" "test_object_1" {
    bucket = "z0-bucket"
    key = "terraform-file-1.txt"
    source = %q
}

resource "qiniu_bucket_object" "test_object_2" {
    bucket = "z0-bucket"
    key = "terraform-file-2.txt"
    source = %q
}

resource "qiniu_bucket_object" "test_object_3" {
    bucket = "z0-bucket"
    key = "terraform-file-3.txt"
    source = %q
}
                `, tmpFiles[0].Name(), tmpFiles[1].Name(), tmpFiles[2].Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_1"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_2"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_3"),
				),
			}, {
				Config: `
data "qiniu_buckets_objects" "all" {
    bucket = "z0-bucket"
    prefix = "terraform-file-"
}

data "qiniu_buckets_objects" "prefixed" {
    bucket = "z0-bucket"
    prefix = "terraform-file-1"
}

data "qiniu_buckets_objects" "limited" {
    bucket = "z0-bucket"
    prefix = "terraform-file-"
    limit = 2
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.all", "keys.#", "3"),
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.all", "key_infos.#", "3"),
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.prefixed", "keys.#", "1"),
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.prefixed", "key_infos.#", "1"),
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.limited", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.qiniu_buckets_objects.limited", "key_infos.#", "2"),
				),
			}},
		})
	})

	It("should verify qiniu buckets objects filter syntax", func() {
		resource.Test(MakeT("TestVerifyQiniuBucketsObjectsFilter"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
data "qiniu_buckets_objects" "no_bucket" {
}
                `,
				ExpectError: regexp.MustCompile("is required, but no definition was found"),
			}, {
				Config: `
data "qiniu_buckets_objects" "no_bucket" {
    bucket = "abc"
    limit = 0
}
                `,
				ExpectError: regexp.MustCompile("must be positive"),
			}, {
				Config: `
data "qiniu_buckets_objects" "no_bucket" {
    bucket = "abc"
    limit = -1
}
                `,
				ExpectError: regexp.MustCompile("must be positive"),
			}},
		})
	})
})
