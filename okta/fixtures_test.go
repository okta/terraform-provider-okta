package okta

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

type fixtureManager struct {
	Path     string
	Seed     int
	TestName string
}

const (
	baseSchema   = "base"
	customSchema = "custom"
	uuidPattern  = "replace_with_uuid"
)

// newFixtureManager Gets a new fixture manager for a particular resource.
func newFixtureManager(resourceName, testName string) *fixtureManager {
	ri := acctest.RandInt()

	// If we are running in VCR mode make the random number be a hash of the
	// test name.
	if os.Getenv("OKTA_VCR_TF_ACC") != "" {
		h := fnv.New32a()
		h.Write([]byte(testName))
		ri = int(h.Sum32())
	}

	dir, _ := os.Getwd()
	return &fixtureManager{
		Path:     path.Join(dir, "../examples", resourceName),
		TestName: testName,
		Seed:     ri,
	}
}

func (manager *fixtureManager) GetFixtures(fixtureName string, t *testing.T) string {
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

	return manager.ConfigReplace(tfConfig)
}

func (manager *fixtureManager) ConfigReplace(tfConfig string) string {
	return strings.ReplaceAll(tfConfig, uuidPattern, fmt.Sprintf("%d", manager.Seed))
}
