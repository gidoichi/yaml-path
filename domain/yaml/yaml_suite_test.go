package yaml_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestYaml(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Yaml Suite")
}
