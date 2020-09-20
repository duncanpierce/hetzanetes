package impl

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"strings"
)

type Errors []error

func (errs Errors) Error() string {
	var errStrings []string
	for _, err := range errs {
		errStrings = append(errStrings, err.Error())
	}
	return strings.Join(errStrings, "\n")
}

func (errs *Errors) Catch(_ *hcloud.Response, err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}
