package plugin_test

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestPlugin(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "plugin suite")
}
