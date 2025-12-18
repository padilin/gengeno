package test

import (
	"os"
	"testing"
)

func init() {
	// When running tests, the CWD is the package directory (test/).
	// We need it to be the project root for assets loading.
	// But we need to be careful not to break if running from root with ./... ?
	// Actually go test ./... runs each pkg in its dir.

	// Check if "assets" exists in current dir.
	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		// Attempt to move up one level
		if _, err := os.Stat("../assets"); err == nil {
			_ = os.Chdir("..")
		}
	}
}

func TestMain(m *testing.M) {
	// This ensures init() runs before tests.
	// Actually init() runs automatically.
	// But TestMain gives us control if needed.
	os.Exit(m.Run())
}
