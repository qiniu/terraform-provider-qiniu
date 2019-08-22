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
			tmpFiles     = make([]*os.File, 3)
			randomString = timeString()
			err          error
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
resource "qiniu_bucket" "basic_bucket" {
    name = "object-test-terraform-%s"
    region_id = "z0"
    private = true
}

resource "qiniu_bucket_object" "test_object_1" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    key = "file-1.txt"
    source = %q
}

resource "qiniu_bucket_object" "test_object_2" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    key = "file-2.txt"
    source = %q
}

resource "qiniu_bucket_object" "test_object_3" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    key = "file-3.txt"
    source = %q
}

data "qiniu_buckets_objects" "all" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
}

data "qiniu_buckets_objects" "prefixed" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    prefix = "test-1"
}

data "qiniu_buckets_objects" "limited" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    limit = 2
}
                `, randomString, tmpFiles[0].Name(), tmpFiles[1].Name(), tmpFiles[2].Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_1"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_2"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object_3"),
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
