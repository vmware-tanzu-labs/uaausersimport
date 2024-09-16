package functions_test

import (
	"os"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	"github.com/pivotalservices/uaausersimport/config"
	"github.com/pivotalservices/uaausersimport/functions"
)

var ctx *config.Context = &config.Context{
	Ccurl:    "https://ccurl.sysdomain.com",
	UAAURL:   "https://uaaurl.sysdomain.com",
	Clientid: "bulkimport",
	Secret:   "test",
}

var _ = Describe("Function", func() {
	file, _ := os.Open("fixtures/users.yml")
	cfg, _ := config.Parse(file)
	logger := lager.NewLogger("uaausersimport")
	ctx.Logger = logger
	ctx.Users = cfg.Users
	var _ = Describe("Map Users", func() {
		It("Should map users with GetTokenFuncs", func() {
			var GetTokenFunc functions.GetTokenFunc = func(*config.Context) (string, error) {
				return "my_token", nil
			}
			Ω(len(ctx.Users)).Should(Equal(2))
			resultFunc := GetTokenFunc.MapUsers()
			userInfos, err := resultFunc(ctx)
			Ω(err).Should(BeNil())
			Ω(len(userInfos)).Should(Equal(2))
			Ω(userInfos[0].Token).Should(Equal("my_token"))
			Ω(reflect.DeepEqual(userInfos[0].User, cfg.Users[0])).Should(Equal(true))
			Ω(reflect.DeepEqual(userInfos[1].User, cfg.Users[1])).Should(Equal(true))
			Ω(userInfos[1].Token).Should(Equal("my_token"))
		})
	})
})
