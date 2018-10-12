package okta

import (
	"fmt"
	"strings"
	"sync"
)

type result struct {
	err error
}

// Dead simple Promise all that will only work in very basic circumstances.
// We could make the function return a generic result object as well.
func promiseAll(wg *sync.WaitGroup, resp chan []*result, funcs ...func() error) {
	resultList := make([]*result, len(funcs))

	for i, f := range funcs {
		wg.Add(1)
		go func(index int, cb func() error) {
			defer wg.Done()
			err := cb()
			resultList[index] = &result{err}
		}(i, f)
	}
	resp <- resultList
}

func getPromiseError(resultList []*result, message string) error {
	var errList []string

	for _, r := range resultList {
		if r.err != nil {
			errList = append(errList, r.err.Error())
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("%s. Errors: %s", message, strings.Join(errList, ", "))
	}

	return nil
}
