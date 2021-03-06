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

var _ = Describe("resourceQiniuBucketObject", func() {
	It("should create qiniu bucket object", func() {
		randomString := timeString()
		tmpFile, err := ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())
		defer os.Remove(tmpFile.Name())
		_, err = io.WriteString(tmpFile, "hello world")
		Expect(err).NotTo(HaveOccurred())
		Expect(tmpFile.Close()).To(Succeed())

		resource.Test(MakeT("TestCreateQiniuBucketObject"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "basic_bucket" {
    name = "terraform-object-test-%s"
    region_id = "z1"
    private = true
}

resource "qiniu_bucket_object" "test_object" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    key = "file-1.txt"
    source = %q
}
                `, randomString, tmpFile.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object"),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "name", fmt.Sprintf("terraform-object-test-%s", randomString)),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "region_id", "z1"),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "private", "true"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "bucket", fmt.Sprintf("terraform-object-test-%s", randomString)),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "key", "file-1.txt"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_length", fmt.Sprintf("%d", len("hello world"))),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_etag", "FiqubDXJT8-0FdvpX0CLnOke6Ebt"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "storage_type", ""),
				),
			}},
		})
	})

	It("should create qiniu bucket object and upload content", func() {
		var (
			randomString = timeString()
			content      = ""
		)
		for i := 0; i < 100; i++ {
			content += "hello world"
		}

		resource.Test(MakeT("TestCreateQiniuBucketObjectByContent"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "basic_bucket" {
    name = "terraform-object-test-%s"
    region_id = "z0"
    private = false
}

resource "qiniu_bucket_object" "test_object" {
    bucket = "${qiniu_bucket.basic_bucket.name}"
    key = "file-2.txt"
    content = %q
    storage_type = "infrequent"
}
	                `, randomString, content),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.basic_bucket"),
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object"),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "name", fmt.Sprintf("terraform-object-test-%s", randomString)),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "region_id", "z0"),
					resource.TestCheckResourceAttr("qiniu_bucket.basic_bucket", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "bucket", fmt.Sprintf("terraform-object-test-%s", randomString)),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "key", "file-2.txt"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_length", fmt.Sprintf("%d", len(content))),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_etag", "FmO9UHw3jb69Wfd4U96mxMLDn37X"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "storage_type", "infrequent"),
				),
			}},
		})
	})

	It("should accept source or content", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucketObject"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket_object" "test_object" {
    bucket = "z0-bucket"
    key = "file.txt"
}
                `,
				ExpectError: regexp.MustCompile("Neither \"source\" nor \"content\" is specified"),
			}},
		})
	})

	It("should accept either source or content", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucketObject"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket_object" "test_object" {
    bucket = "z0-bucket"
    key = "file.txt"
    source = "/etc/services"
    content = "abcdef"
}
                `,
				ExpectError: regexp.MustCompile("conflicts with"),
			}},
		})
	})

	It("should reject if source is invalid", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucketObject"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket_object" "test_object" {
    bucket = "z0-bucket"
    key = "file.txt"
    source = "/not/existed"
}
                `,
				ExpectError: regexp.MustCompile("no such file or directory"),
			}},
		})
	})

	It("should reject if storage_type is invalid", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucketObject"), resource.TestCase{
			PreCheck:  testPreCheck,
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket_object" "test_object" {
    bucket = "z0-bucket"
    key = "file.txt"
    content = "abcdef"
    storage_type = "invalid_type"
}
                `,
				ExpectError: regexp.MustCompile("invalid object storage type"),
			}},
		})
	})
})
