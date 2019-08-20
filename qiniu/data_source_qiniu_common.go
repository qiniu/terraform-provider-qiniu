package qiniu

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/hashcode"
)

func dataResourceIdHash(ids []string) string {
	var buf bytes.Buffer

	for _, id := range ids {
		buf.WriteString(fmt.Sprintf("%s-", id))
	}

	return strconv.Itoa(hashcode.String(buf.String()))
}
