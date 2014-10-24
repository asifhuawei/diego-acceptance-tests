package diego

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/diego-acceptance-tests/helpers/assets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Getting instance information", func() {
	var appName string

	BeforeEach(func() {
		appName = generator.RandomName()

		Eventually(cf.Cf("push", appName, "-p", assets.NewAssets().Standalone, "--no-start", "-b", DIEGO_NULL_BUILDPACK), CF_PUSH_TIMEOUT).Should(Exit(0))
		Eventually(cf.Cf("set-env", appName, "DIEGO_STAGE_BETA", "true")).Should(Exit(0))
		Eventually(cf.Cf("set-env", appName, "DIEGO_RUN_BETA", "true")).Should(Exit(0))

		Eventually(cf.Cf("scale", appName, "-i", "3")).Should(Exit(0))
		Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(cf.Cf("delete", appName, "-f")).Should(Exit(0))
	})

	It("Retrieves instance information for cf app and cf apps", func() {
		By("calling cf app")
		app := cf.Cf("app", appName).Wait()
		Expect(app).To(Exit(0))
		Expect(app).To(Say("instances: [0-3]/3"))
		Expect(app).To(Say("#0"))
		Expect(app).To(Say("#1"))
		Expect(app).To(Say("#2"))
		Expect(app).ToNot(Say("#3"))

		By("calling cf apps")
		app = cf.Cf("apps").Wait()
		Expect(app).To(Exit(0))
		Expect(app).To(Say(appName))
		Expect(app).To(Say("[0-3]/3"))
	})
})
