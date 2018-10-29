package okta

import (
	"fmt"
	"strings"
	"sync"
)

// Basic result struct that can be expanded to support output
type result struct {
	err error
}

// Dead simple Promise all that will only work in very basic circumstances.
// We could make the function return a generic result object as well.
func promiseAll(limit int, wg *sync.WaitGroup, resp chan []*result, funcs ...func() error) {
	jobIndex := 0
	resultList := make([]*result, len(funcs))

	for jobIndex < len(funcs) {
		for i := 0; i < limit; i++ {
			wg.Add(1)
			go func(index int, cb func() error) {
				defer wg.Done()
				err := cb()
				resultList[index] = &result{err}
			}(jobIndex, funcs[jobIndex])
			jobIndex++
		}
		wg.Wait()
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
