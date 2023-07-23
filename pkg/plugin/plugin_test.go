package plugin_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/containernetworking/cni/pkg/skel"

	"github.com/eddev/template-cni/pkg/plugin"
)

var _ = ginkgo.Describe("template-cni", func() {
	ginkgo.It("Add", func() {
		args := &skel.CmdArgs{
			StdinData: []byte("{}"),
		}
		gomega.Expect(plugin.CmdAdd(args)).To(gomega.Succeed())
	})

	ginkgo.It("Del", func() {
		gomega.Expect(plugin.CmdDel(&skel.CmdArgs{})).To(gomega.Succeed())
	})

	ginkgo.It("Check", func() {
		gomega.Expect(plugin.CmdCheck(&skel.CmdArgs{})).To(gomega.Succeed())
	})
})
