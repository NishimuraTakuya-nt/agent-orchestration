package use_cases_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUseCases(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UseCases Suite")
}