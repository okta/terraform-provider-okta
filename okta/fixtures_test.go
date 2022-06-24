package okta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

type fixtureManager struct {
	Path string
}

const (
	baseSchema   = "base"
	customSchema = "custom"
	uuidPattern  = "replace_with_uuid"
)

// newFixtureManager get a new fixture manager for a particular resource.
func newFixtureManager(resourceName string) *fixtureManager {
	dir, _ := os.Getwd()
	return &fixtureManager{
		Path: path.Join(dir, "../examples", resourceName),
	}
}

func (manager *fixtureManager) GetFixtures(fixtureName string, rInt int, t *testing.T) string {
	file, err := os.Open(path.Join(manager.Path, fixtureName))
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}
	defer file.Close()
	rawFile, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to load terraform fixtures for ACC test, err: %v", err)
	}
	tfConfig := string(rawFile)
	if strings.Count(tfConfig, uuidPattern) == 0 {
		return tfConfig
	}

	return manager.ConfigReplace(tfConfig, rInt)
}

func (manager *fixtureManager) ConfigReplace(tfConfig string, rInt int) string {
	return strings.ReplaceAll(tfConfig, uuidPattern, fmt.Sprintf("%d", rInt))
}
