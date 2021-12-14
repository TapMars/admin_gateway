package config_test

import (
	"TapMars/productManager/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Config", func() {
	Describe("Get Environment Variables", func() {
		BeforeEach(func() {
			_ = os.Setenv("PORT", "")
			_ = os.Setenv("PROJECT_ID", "")
		})

		It("Should Pass with both values populated", func() {
			_ = os.Setenv("PORT", "8080")
			_ = os.Setenv("PROJECT_ID", "happy")
			port, projectID, err := config.GetEnvironmentVariables()

			Expect(port).To(Equal("8080"))
			Expect(projectID).To(Equal("happy"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should Fail with Port value populated", func() {
			_ = os.Setenv("PORT", "8080")
			_, _, err := config.GetEnvironmentVariables()

			Expect(err).To(HaveOccurred())
		})

		It("Should Fail with Project_id value populated", func() {
			_ = os.Setenv("PROJECT_ID", "happy")
			_, _, err := config.GetEnvironmentVariables()

			Expect(err).To(HaveOccurred())
		})
	})

})
