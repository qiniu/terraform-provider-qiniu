package qiniu_test

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("resourceQiniuBucket", func() {
	It("should create qiniu bucket", func() {
		resourceID := "qiniu_bucket.basic_bucket"
		resource.Test(MakeT("TestCreateQiniuBucket"), resource.TestCase{
			PreCheck:      testPreCheck,
			IDRefreshName: resourceID,
			Providers:     providers,
			CheckDestroy:  testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "basic_bucket" {
    name = "basic-test-terraform"
    region_id = "z2"
    private = true
    index_page_on = false
    lifecycle_rules {
        name = "rule_for_user_files"
        prefix = "users/"
        to_line_after_days = 30
    }
    lifecycle_rules {
        name = "rule_for_sys_files"
        prefix = "sys/"
        delete_after_days = 10
    }
    cors_rules {
        allowed_origins = ["http://www.qiniu.com"]
        allowed_methods = ["GET", "POST"]
    }
    anti_leech_mode = "whitelist"
    referer_pattern = "*.qiniu.com;*.qiniudn.com"
    allow_empty_referer = true
    only_enable_anti_leech_for_cdn = false
    max_age = 86400
    tagging = {
        env = "test"
        kind = "basic"
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "basic-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "true"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceID, "cors_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceID, "anti_leech_mode", "whitelist"),
					resource.TestCheckResourceAttr(resourceID, "referer_pattern", "*.qiniu.com;*.qiniudn.com"),
					resource.TestCheckResourceAttr(resourceID, "allow_empty_referer", "true"),
					resource.TestCheckResourceAttr(resourceID, "only_enable_anti_leech_for_cdn", "false"),
					resource.TestCheckResourceAttr(resourceID, "max_age", "86400"),
					resource.TestCheckResourceAttr(resourceID, "tagging.%", "2"),
					resource.TestCheckResourceAttr(resourceID, "tagging.env", "test"),
					resource.TestCheckResourceAttr(resourceID, "tagging.kind", "basic"),
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
}
                `,
				ExpectError: regexp.MustCompile("must not contain invalid characters"),
			}},
		})
	})

	It("should reject empty qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = ""
    region_id = "z2"
}
                `,
				ExpectError: regexp.MustCompile("must not be empty"),
			}},
		})
	})

	It("should reject too short qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "invalid_bucket" {
    name = "ab"
    region_id = "z2"
}
                `,
				ExpectError: regexp.MustCompile("must not be shorter than 3 characters"),
			}},
		})
	})

	It("should reject too long qiniu bucket name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "longlonglonglonglonglonglonglonglonglonglonglonglonglonglonglong-%s"
    region_id = "z2"
}
                `, timeString()),
				ExpectError: regexp.MustCompile("must not be longer than 63 characters"),
			}},
		})
	})

	It("should reject invalid qiniu region id", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z100"
}
                `, timeString()),
				ExpectError: regexp.MustCompile("is invalid"),
			}},
		})
	})

	It("should reject invalid qiniu image_url", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    image_url = "www.qiniu.com"
}
                `, timeString()),
				ExpectError: regexp.MustCompile("must be valid url"),
			}},
		})
	})

	It("should reject invalid qiniu image_host", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    image_url = "http://www.qiniu.com"
    image_host = "http://www.qiniu.com"
}
                `, timeString()),
				ExpectError: regexp.MustCompile("must be valid host"),
			}},
		})
	})

	It("should reject invalid bucket lifecycle rule name", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    lifecycle_rules {
        name = "superlongsuperlongsuperlongsuperlongsuperlongsuperlong"
    }
}
                `, timeString()),
				ExpectError: regexp.MustCompile("must not be longer than and equal to 50 characters"),
			}},
		})
	})

	It("should reject invalid bucket anti leech mode", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    anti_leech_mode = "invalid"
}
                `, timeString()),
				ExpectError: regexp.MustCompile("\"anti_leech_mode\" contains invalid mode"),
			}},
		})
	})

	It("should reject invalid cors rule without allowed origins", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    cors_rules {
        allowed_methods = ["GET"]
    }
}
                `, timeString()),
				ExpectError: regexp.MustCompile("The argument \"allowed_origins\" is required, but no definition was found."),
			}},
		})
	})

	It("should reject invalid cors rule without allowed methods", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    cors_rules {
        allowed_origins = ["http://abc.com"]
    }
}
                `, timeString()),
				ExpectError: regexp.MustCompile("The argument \"allowed_methods\" is required, but no definition was found."),
			}},
		})
	})

	It("should reject invalid cors rule with empty allowed origins", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    cors_rules {
        allowed_origins = []
        allowed_methods = ["GET"]
    }
}
                `, timeString()),
				ExpectError: regexp.MustCompile("invalid argument"),
			}},
		})
	})

	It("should reject invalid cors rule with invalid method", func() {
		resource.Test(MakeT("TestCreateInvalidQiniuBucket"), resource.TestCase{
			Providers: providers,
			Steps: []resource.TestStep{{
				Config: fmt.Sprintf(`
resource "qiniu_bucket" "invalid_bucket" {
    name = "valid-name-%s"
    region_id = "z1"
    cors_rules {
        allowed_origins = ["http://abc.com"]
        allowed_methods = ["OPEN"]
    }
}
                `, timeString()),
				ExpectError: regexp.MustCompile("invalid http method"),
			}},
		})
	})

	It("should update qiniu bucket", func() {
		resourceID := "qiniu_bucket.update_bucket"
		resource.Test(MakeT("TestUpdateQiniuBucket"), resource.TestCase{
			PreCheck:      testPreCheck,
			IDRefreshName: resourceID,
			Providers:     providers,
			CheckDestroy:  testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "max_age", "0"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    max_age = "86400"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "max_age", "86400"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    max_age = "172800"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "max_age", "172800"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
    image_url = "http://www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "max_age", "0"),
					resource.TestCheckResourceAttr(resourceID, "image_url", "http://www.qiniu.com"),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = false
    image_url = "http://portal.qiniu.io"
    image_host = "www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", "http://portal.qiniu.io"),
					resource.TestCheckResourceAttr(resourceID, "image_host", "www.qiniu.com"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "true"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    index_page_on = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "true"),
					resource.TestCheckResourceAttr(resourceID, "image_url", ""),
					resource.TestCheckResourceAttr(resourceID, "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    lifecycle_rules {
        name = "rule_for_user_files"
        prefix = "users/"
        to_line_after_days = 30
    }
    lifecycle_rules {
        name = "rule_for_sys_files"
        prefix = "sys/"
        delete_after_days = 10
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "2"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    lifecycle_rules {
        name = "rule_for_admin_files"
        prefix = "admins/"
        delete_after_days = 20
    }
    lifecycle_rules {
        name = "rule_for_sys_files"
        prefix = "sys/"
        to_line_after_days = 50
    }
    lifecycle_rules {
        name = "rule_for_guest_files"
        prefix = "guests/"
        delete_after_days = 30
        to_line_after_days = 10
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "3"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    anti_leech_mode = "whitelist"
    referer_pattern = "*.qiniu.com;*.qiniudn.com"
    allow_empty_referer = true
    only_enable_anti_leech_for_cdn = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceID, "anti_leech_mode", "whitelist"),
					resource.TestCheckResourceAttr(resourceID, "referer_pattern", "*.qiniu.com;*.qiniudn.com"),
					resource.TestCheckResourceAttr(resourceID, "allow_empty_referer", "true"),
					resource.TestCheckResourceAttr(resourceID, "only_enable_anti_leech_for_cdn", "true"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    anti_leech_mode = "blacklist"
    referer_pattern = "*.qiniu.com;*.qiniudn.com"
    allow_empty_referer = false
    only_enable_anti_leech_for_cdn = false
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceID, "anti_leech_mode", "blacklist"),
					resource.TestCheckResourceAttr(resourceID, "referer_pattern", "*.qiniu.com;*.qiniudn.com"),
					resource.TestCheckResourceAttr(resourceID, "allow_empty_referer", "false"),
					resource.TestCheckResourceAttr(resourceID, "only_enable_anti_leech_for_cdn", "false"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "lifecycle_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceID, "anti_leech_mode", ""),
					resource.TestCheckResourceAttr(resourceID, "referer_pattern", ""),
					resource.TestCheckResourceAttr(resourceID, "allow_empty_referer", "false"),
					resource.TestCheckResourceAttr(resourceID, "only_enable_anti_leech_for_cdn", "false"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    cors_rules {
        allowed_origins = ["http://*.abc.com", "http://*.def.com"]
        allowed_methods = ["GET"]
        allowed_headers = ["Content-Type"]
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "cors_rules.#", "1"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    cors_rules {
        allowed_origins = ["http://*.abc.com"]
        allowed_methods = ["GET"]
        allowed_headers = ["Content-Type"]
    }
    cors_rules {
        allowed_origins = ["http://*.def.com"]
        allowed_methods = ["POST"]
        allowed_headers = ["Content-Type", "Content-Encoding"]
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "cors_rules.#", "2"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    tagging = {
        environment = "test"
        kind = "basic"
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "cors_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceID, "tagging.%", "2"),
					resource.TestCheckResourceAttr(resourceID, "tagging.environment", "test"),
					resource.TestCheckResourceAttr(resourceID, "tagging.kind", "basic"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
    tagging = {
        environment = "production"
        kind = "advanced"
        user = "bachue"
    }
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "tagging.%", "3"),
					resource.TestCheckResourceAttr(resourceID, "tagging.environment", "production"),
					resource.TestCheckResourceAttr(resourceID, "tagging.kind", "advanced"),
					resource.TestCheckResourceAttr(resourceID, "tagging.user", "bachue"),
				),
			}, {
				Config: `
resource "qiniu_bucket" "update_bucket" {
    name = "update-test-terraform"
    region_id = "z2"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists(resourceID),
					resource.TestCheckResourceAttr(resourceID, "name", "update-test-terraform"),
					resource.TestCheckResourceAttr(resourceID, "region_id", "z2"),
					resource.TestCheckResourceAttr(resourceID, "private", "false"),
					resource.TestCheckResourceAttr(resourceID, "index_page_on", "false"),
					resource.TestCheckResourceAttr(resourceID, "tagging.%", "0"),
				),
			}},
		})
	})

	It("should create qiniu buckets", func() {
		resource.Test(MakeT("TestCreateQiniuBuckets"), resource.TestCase{
			PreCheck:     testPreCheck,
			Providers:    providers,
			CheckDestroy: testCheckQiniuResourceDestroy,
			Steps: []resource.TestStep{{
				Config: `
resource "qiniu_bucket" "public-bucket" {
    name = "bucket-test-1-terraform"
    region_id = "z2"
}

resource "qiniu_bucket" "private-bucket" {
    name = "bucket-test-2-terraform"
    region_id = "z2"
    private = true
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.public-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "name", "bucket-test-1-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "image_url", ""),
					resource.TestCheckResourceAttr("qiniu_bucket.public-bucket", "image_host", ""),
					testCheckQiniuBucketItemExists("qiniu_bucket.private-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "name", "bucket-test-2-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "private", "true"),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "image_url", ""),
					resource.TestCheckResourceAttr("qiniu_bucket.private-bucket", "image_host", ""),
				),
			}, {
				Config: `
resource "qiniu_bucket" "imaged-bucket" {
    name = "bucket-test-3-terraform"
    region_id = "z2"
    image_url = "http://qiniu.io"
}

resource "qiniu_bucket" "imaged-bucket-with-host" {
    name = "bucket-test-4-terraform"
    region_id = "z2"
    image_url = "http://qiniu.io"
    image_host = "www.qiniu.com"
}
                `,
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckQiniuBucketItemExists("qiniu_bucket.imaged-bucket"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "name", "bucket-test-3-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "image_url", "http://qiniu.io"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket", "image_host", ""),
					testCheckQiniuBucketItemExists("qiniu_bucket.imaged-bucket-with-host"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "name", "bucket-test-4-terraform"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "region_id", "z2"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "private", "false"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "image_url", "http://qiniu.io"),
					resource.TestCheckResourceAttr("qiniu_bucket.imaged-bucket-with-host", "image_host", "www.qiniu.com"),
				),
			}},
		})
	})
})
