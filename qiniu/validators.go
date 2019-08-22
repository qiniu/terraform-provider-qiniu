package qiniu

import (
	"fmt"
	"regexp"
)

var validBucketNameRegex *regexp.Regexp

func init() {
	validBucketNameRegex = regexp.MustCompile("^[a-zA-Z0-9\\-_]+$")
}

func validateBucketName(v interface{}, attributeName string) (warns []string, errs []error) {
	bucketName := v.(string)
	if len(bucketName) == 0 {
		errs = append(errs, fmt.Errorf("%q must not be empty", attributeName))
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

func validateRegex(v interface{}, attributeName string) (warns []string, errs []error) {
	if _, err := regexp.Compile(v.(string)); err != nil {
		errs = append(errs, fmt.Errorf("%q contains an invalid regular expression", attributeName))
	}
	return
}

func validatePositiveInt(v interface{}, attributeName string) (warns []string, errs []error) {
	if v.(int) <= 0 {
		errs = append(errs, fmt.Errorf("%q must be positive", attributeName))
	}
	return
}
