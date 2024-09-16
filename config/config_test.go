package config_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/uaausersimport/config"
)

var _ = Describe("Config", func() {
	file, _ := os.Open("fixtures/users.yml")
	Describe("Parse config", func() {
		Context("with correct fixture", func() {
			It("should be return correct config", func() {
				cfg, err := config.Parse(file)
				Ω(err).Should(BeNil())
				Ω(len(cfg.Users)).Should(Equal(2))
				Ω(cfg.Users[0].Externalid).Should(Equal("uid=jcalabrese,ou=People,dc=homelab,dc=io"))
				Ω(cfg.Users[1].Externalid).Should(Equal("uid=rparrish,ou=People,dc=homelab,dc=io"))
				Ω(cfg.Users[0].Emails[0]).Should(Equal("jcalabrese@pivotal.io"))
				Ω(cfg.Users[1].Emails[0]).Should(Equal("rparrish@pivotal.io"))
				Ω(len(cfg.Users[0].Orgs)).Should(Equal(2))
				Ω(cfg.Users[0].Orgs[0].Roles[0]).Should(Equal("OrgManager"))
				Ω(cfg.Users[0].Orgs[0].Roles[1]).Should(Equal("OrgAuditor"))
				Ω(len(cfg.Users[0].Orgs[0].Spaces)).Should(Equal(2))
				Ω(cfg.Users[0].Orgs[0].Spaces[0].Roles[0]).Should(Equal("SpaceManager"))
				Ω(cfg.Users[0].Orgs[0].Spaces[1].Roles[1]).Should(Equal("SpaceAuditor"))
				Ω(cfg.Origin).Should(Equal("ldap"))
				Ω(len(cfg.Users[0].Orgs[1].Spaces)).Should(Equal(2))
			})
		})
	})
})
