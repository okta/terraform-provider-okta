package okta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

type (
	fixtureManager struct {
		Path string
	}
)

// newFixtureManager get a new fixture manager for a particular resource.
func newFixtureManager(resourceName string) *fixtureManager {
	_, filename, _, _ := runtime.Caller(0)
	exPath := filepath.Dir(filename)
	return &fixtureManager{
		Path: path.Join(exPath, "../examples", resourceName),
	}
}

func (manager *fixtureManager) GetFixtures(fixtureName string, rInt int, t *testing.T) string {
	file, err := os.Open(path.Join(manager.Path, fixtureName))
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}
	rawFile, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}

	// Template string
	tmpString := fmt.Sprintf(string(rawFile), rInt)
	return tmpString
}
