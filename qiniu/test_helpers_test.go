package qiniu_test

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	qiniu "github.com/qiniu/terraform-provider-qiniu/qiniu"
)

var (
	qiniuProvider *schema.Provider
	providers     map[string]terraform.ResourceProvider
)

func init() {
	qiniuProvider = qiniu.Provider().(*schema.Provider)
	providers = map[string]terraform.ResourceProvider{
		"qiniu": qiniuProvider,
	}
}

func testPreCheck() {
	Expect(os.Getenv("QINIU_ACCESS_KEY")).NotTo(BeEmpty())
	Expect(os.Getenv("QINIU_SECRET_KEY")).NotTo(BeEmpty())
}

func timeString() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

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

func testCheckQiniuBucketObjectItemExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client := qiniuProvider.Meta().(*qiniu.Client)
		bucketName, key := getBucketNameAndKeyFromEntry(rs.Primary.ID)
		_, err := client.BucketManager.Stat(bucketName, key)
		return err
	}
}

func testCheckQiniuResourceDestroy(s *terraform.State) (err error) {
	client := qiniuProvider.Meta().(*qiniu.Client)

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "qiniu_bucket":
			bucketName := rs.Primary.ID
			if _, err = client.BucketManager.GetBucketInfo(bucketName); err == nil {
				return fmt.Errorf("Bucket still exists")
			} else if !qiniu.IsResourceNotFound(err) {
				return
			}
		case "qiniu_bucket_object":
			bucketName, key := getBucketNameAndKeyFromEntry(rs.Primary.ID)
			if _, err = client.BucketManager.Stat(bucketName, key); err == nil {
				return fmt.Errorf("Bucket Object still exists")
			} else if !qiniu.IsResourceNotFound(err) {
				return
			}
		}
	}

	return nil
}

type T struct {
	ginkgoT GinkgoTInterface
	name    string
}

func MakeT(name string) resource.TestT {
	return &T{ginkgoT: GinkgoT(), name: name}
}

func (t *T) Error(args ...interface{}) {
	t.ginkgoT.Error(args...)
}

func (t *T) Fatal(args ...interface{}) {
	t.ginkgoT.Fatal(args...)
}

func (t *T) Skip(args ...interface{}) {
	t.ginkgoT.Skip(args...)
}

func (t *T) Name() string {
	return t.name
}

func (t *T) Parallel() {
	t.ginkgoT.Parallel()
}

func getBucketNameAndKeyFromEntry(entry string) (bucket string, key string) {
	parts := strings.SplitN(entry, ":", 2)
	return parts[0], parts[1]
}

func getEntryFromBucketNameAndKey(bucket string, key string) (entry string) {
	return bucket + ":" + key
}
