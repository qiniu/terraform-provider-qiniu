package qiniu_test

import (
	"github.com/hashicorp/terraform/helper/schema"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/qiniu/terraform-provider-qiniu/qiniu"
)

var _ = Describe("Provider", func() {
	It("should pass internal validation", func() {
		Expect(qiniu.Provider().(*schema.Provider).InternalValidate()).To(Succeed())
	})
})
