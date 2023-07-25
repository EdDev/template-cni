package plugin_test

import (
	"bytes"
	"fmt"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/containernetworking/plugins/pkg/testutils"

	"github.com/eddev/template-cni/pkg/plugin"
)

var _ = ginkgo.Describe("template-cni", func() {
	var testNS ns.NetNS

	ginkgo.BeforeEach(func() {
		var err error
		testNS, err = testutils.NewNS()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		ginkgo.DeferCleanup(testNS.Close)
	})

	ginkgo.It("Add", func() {
		args := &skel.CmdArgs{
			ContainerID: "123456789",
			Netns:       testNS.Path(),
			IfName:      "dummy0",
			StdinData:   []byte(`{"cniVersion":"1.0.0"}`),
		}
		result, err := plugin.CmdAddResult(args)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		versionedResult, err := result.GetAsVersion(result.Version())
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		var buf bytes.Buffer
		gomega.Expect(versionedResult.PrintTo(&buf)).To(gomega.Succeed())
		gomega.Expect(buf.String()).To(gomega.MatchJSON(fmt.Sprintf(`
			{
				"cniVersion": "1.0.0",
				"interfaces": [
					{
						"name": %q,
						"mac": "00:00:00:00:00:01",
						"sandbox": %q
					}
				],
				"dns": {}
			}
		`, args.IfName, testNS.Path())))
	})

	ginkgo.It("Del", func() {
		gomega.Expect(plugin.CmdDel(&skel.CmdArgs{})).To(gomega.Succeed())
	})

	ginkgo.It("Check", func() {
		gomega.Expect(plugin.CmdCheck(&skel.CmdArgs{})).To(gomega.Succeed())
	})
})
