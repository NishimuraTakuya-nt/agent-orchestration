package testutils

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// SetupGinkgoConfig configures Ginkgo for testing
func SetupGinkgoConfig() {
	// Set default timeout for Eventually/Consistently
	SetDefaultEventuallyTimeout(10)
	SetDefaultEventuallyPollingInterval(100)
	SetDefaultConsistentlyDuration(1)
	SetDefaultConsistentlyPollingInterval(100)
	
	// Configure fail handler
	RegisterFailHandler(Fail)
}

// IsIntegrationTest checks if integration tests should run
func IsIntegrationTest() bool {
	return os.Getenv("INTEGRATION_TEST") == "true"
}

// IsE2ETest checks if E2E tests should run
func IsE2ETest() bool {
	return os.Getenv("E2E_TEST") == "true"
}

// SkipIfNotIntegration skips test if not running integration tests
func SkipIfNotIntegration() {
	if !IsIntegrationTest() {
		Skip("Skipping integration test - set INTEGRATION_TEST=true to run")
	}
}

// SkipIfNotE2E skips test if not running E2E tests
func SkipIfNotE2E() {
	if !IsE2ETest() {
		Skip("Skipping E2E test - set E2E_TEST=true to run")
	}
}