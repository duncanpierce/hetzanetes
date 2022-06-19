package catch

import (
	"github.com/hetznercloud/hcloud-go/hcloud"
	"strings"
	"time"
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
	errs.Add(err)
}

func (errs Errors) OrNil() error {
	if len(errs) == 0 {
		return nil
	}
	return errs
}

func (errs *Errors) Add(err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

func (errs *Errors) HasErrors() bool {
	return len(*errs) > 0
}

func (errs *Errors) Retry(times int, sleep time.Duration, action func() error) {
	err := action()
	for i := 1; i < times; i++ {
		if err == nil {
			return
		}
		time.Sleep(sleep)
		err = action()
	}
	errs.Add(err)
}
