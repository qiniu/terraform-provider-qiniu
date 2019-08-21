package qiniu_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("resourceQiniuBucketObject", func() {
	It("should create qiniu bucket object", func() {
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
resource "qiniu_bucket_object" "test_object" {
    bucket = "z1-bucket"
    key = "file-1.txt"
    source = %q
}
                `, tmpFile.Name()),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "bucket", "z1-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "key", "file-1.txt"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_length", fmt.Sprintf("%d", len("hello world"))),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_etag", "FiqubDXJT8-0FdvpX0CLnOke6Ebt"),
				),
			}},
		})
	})

	It("should create qiniu bucket object and upload content", func() {
		var content = ""
		for i := 0; i < 100; i++ {
			content += "hello world"
		}

		resource.Test(MakeT("TestCreateQiniuBucketObjectByContent"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket_object" "test_object" {
    bucket = "z2-bucket"
    key = "file-2.txt"
    content = %q
}
	                `, content),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketObjectItemExists("qiniu_bucket_object.test_object"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "bucket", "z2-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "key", "file-2.txt"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_length", fmt.Sprintf("%d", len(content))),
					resource.TestCheckResourceAttr("qiniu_bucket_object.test_object", "content_etag", "FmO9UHw3jb69Wfd4U96mxMLDn37X"),
				),
			}},
		})
	})
})
