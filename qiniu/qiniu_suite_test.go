package qiniu_test

import (
	"testing"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQiniu(t *testing.T) {
	BeforeSuite(func() {
		_ = godotenv.Load("../.env")
	})

	RegisterFailHandler(Fail)
	RunSpecs(t, "Qiniu Suite")
}
