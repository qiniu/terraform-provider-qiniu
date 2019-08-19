package qiniu_test

import (
	"os"

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
