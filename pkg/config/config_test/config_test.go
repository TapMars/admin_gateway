package config_test

import (
	"TapMars/admin_gateway/pkg/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {
	Describe("Get Environment Variables", func() {
		BeforeEach(func() {
			_ = os.Setenv("PORT", "")
			_ = os.Setenv("HOST", "")
		})

		It("Should Pass with both values populated", func() {
			_ = os.Setenv("PORT", "8080")
			_ = os.Setenv("HOST", "happy")
			port, host, err := config.GetEnvironmentVariables()

			Expect(port).To(Equal("8080"))
			Expect(host).To(Equal("happy"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should Fail with Port value populated", func() {
			_ = os.Setenv("PORT", "8080")
			_, _, err := config.GetEnvironmentVariables()

			Expect(err).To(HaveOccurred())
		})

		It("Should Fail with Project_id value populated", func() {
			_ = os.Setenv("HOST", "happy")
			_, _, err := config.GetEnvironmentVariables()

			Expect(err).To(HaveOccurred())
		})
	})

})
