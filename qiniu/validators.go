package qiniu

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	validBucketNameRegex *regexp.Regexp
	validRuleNameRegex   *regexp.Regexp
)

func init() {
	validBucketNameRegex = regexp.MustCompile("^[a-zA-Z0-9\\-]+$")
	validRuleNameRegex = regexp.MustCompile("^[a-zA-Z0-9\\_]+$")
}

func validateBucketName(v interface{}, attributeName string) (warns []string, errs []error) {
	bucketName := v.(string)
	if len(bucketName) == 0 {
		errs = append(errs, fmt.Errorf("%q must not be empty", attributeName))
		return
	}
	if len(bucketName) < 3 {
		errs = append(errs, fmt.Errorf("%q must not be shorter than 3 characters", attributeName))
		return
	}
	if len(bucketName) > 63 {
		errs = append(errs, fmt.Errorf("%q must not be longer than 63 characters", attributeName))
		return
	}
	if !validBucketNameRegex.MatchString(bucketName) {
		errs = append(errs, fmt.Errorf("%q must not contain invalid characters", attributeName))
		return
	}
	return
}

func validateRegionID(v interface{}, attributeName string) (warns []string, errs []error) {
	regionId := v.(string)
	switch regionId {
	case "z0", "z1", "z2", "na0", "as0":
		return
	default:
		errs = append(errs, fmt.Errorf("%q is invalid", attributeName))
		return
	}
}

func validatePositiveInt(v interface{}, attributeName string) (warns []string, errs []error) {
	if v.(int) <= 0 {
		errs = append(errs, fmt.Errorf("%q must be positive", attributeName))
	}
	return
}

func validateURL(v interface{}, attributeName string) (warns []string, errs []error) {
	urlString := v.(string)
	if urlString == "" {
		return
	}
	if u, err := url.ParseRequestURI(urlString); err != nil {
		errs = append(errs, fmt.Errorf("%q must be valid url", attributeName))
	} else if u.Scheme != "http" && u.Scheme != "https" {
		errs = append(errs, fmt.Errorf("%q should be http or https protocol", attributeName))
	}
	return
}

func validateHost(v interface{}, attributeName string) (warns []string, errs []error) {
	const r = "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$"
	hostString := v.(string)
	if hostString == "" {
		return
	}
	if !regexp.MustCompile(r).MatchString(hostString) {
		errs = append(errs, fmt.Errorf("%q must be valid host", attributeName))
	}
	return
}

func validateLifecycleRuleName(v interface{}, attributeName string) (warns []string, errs []error) {
	ruleName := v.(string)
	if len(ruleName) == 0 {
		errs = append(errs, fmt.Errorf("%q must not be empty", attributeName))
		return
	}
	if len(ruleName) >= 50 {
		errs = append(errs, fmt.Errorf("%q must not be longer than and equal to 50 characters", attributeName))
		return
	}
	if !validRuleNameRegex.MatchString(ruleName) {
		errs = append(errs, fmt.Errorf("%q must not contain invalid characters", attributeName))
		return
	}
	return
}

func validateAntiLeechMode(v interface{}, attributeName string) (warns []string, errs []error) {
	modeName := v.(string)
	switch strings.ToLower(modeName) {
	case "":
	case "whitelist":
	case "blacklist":
	default:
		errs = append(errs, fmt.Errorf("%q contains invalid mode", attributeName))
	}
	return
}

func validateHTTPMethods(v interface{}, attributeName string) (warns []string, errs []error) {
	switch strings.ToLower(v.(string)) {
	case "get":
	case "head":
	case "post":
	case "put":
	case "delete":
	case "patch":
	case "options":
	case "connect":
	case "trace":
	default:
		errs = append(errs, fmt.Errorf("%q is an invalid http method", attributeName))
	}
	return
}

func validateObjectStorageType(v interface{}, attributeName string) (warns []string, errs []error) {
	switch strings.ToLower(v.(string)) {
	case NormalStorage:
	case InfrequentStorage:
	default:
		errs = append(errs, fmt.Errorf("%q is an invalid object storage type", attributeName))
	}
	return
}
