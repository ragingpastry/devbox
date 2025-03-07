package impl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/stretchr/testify/assert"

	"go.jetpack.io/devbox/internal/env"
	"go.jetpack.io/devbox/internal/fileutil"
	"go.jetpack.io/devbox/internal/planner/plansdk"
)

func TestDevbox(t *testing.T) {
	t.Setenv("TMPDIR", "/tmp")
	t.Setenv(env.DevboxDoNotUpgradeConfig, "1")
	testPaths, err := doublestar.FilepathGlob("../../examples/**/devbox.json")
	assert.NoError(t, err, "Reading testdata/ should not fail")

	assert.Greater(t, len(testPaths), 0, "testdata/ and examples/ should contain at least 1 test")

	for _, testPath := range testPaths {
		if !strings.Contains(testPath, "/commands/") {
			testShellPlan(t, testPath)
		}
	}
}
func testShellPlan(t *testing.T, testPath string) {
	baseDir := filepath.Dir(testPath)
	testName := fmt.Sprintf("%s_shell_plan", filepath.Base(baseDir))
	t.Run(testName, func(t *testing.T) {
		t.Setenv(env.XDGDataHome, "/tmp/devbox")
		assert := assert.New(t)
		shellPlanFile := filepath.Join(baseDir, "shell_plan.json")
		hasShellPlanFile := fileutil.Exists(shellPlanFile)

		box, err := Open(baseDir, os.Stdout)
		assert.NoErrorf(err, "%s should be a valid devbox project", baseDir)

		shellPlan, err := box.ShellPlan()
		assert.NoError(err, "devbox shell plan should not fail")

		if hasShellPlanFile {
			data, err := os.ReadFile(shellPlanFile)
			assert.NoError(err, "shell_plan.json should be readable")

			expected := &plansdk.ShellPlan{}
			err = json.Unmarshal(data, &expected)
			assert.NoError(err, "plan.json should parse correctly")
			assertShellPlansMatch(t, expected, shellPlan)
		}
	})
}

func assertShellPlansMatch(t *testing.T, expected *plansdk.ShellPlan, actual *plansdk.ShellPlan) {
	assert := assert.New(t)

	assert.ElementsMatch(expected.DevPackages, actual.DevPackages, "DevPackages should match")
}
