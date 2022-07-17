package searcher_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSearcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Searcher Suite")
}
