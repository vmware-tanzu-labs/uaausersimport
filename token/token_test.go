package token_test

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
	. "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/gtils/http/httptest"
	"github.com/pivotalservices/uaausersimport/config"
	. "github.com/pivotalservices/uaausersimport/token"
)

var requestAdaptor RequestAdaptor = func() (*http.Response, error) {
	file, _ := os.Open("fixtures/response.json")
	fakeResponse := http.Response{
		Body: file,
	}
	return &fakeResponse, nil
}

var mockGateway *httptest.MockGateway = &httptest.MockGateway{
	FakePostAdaptor: requestAdaptor,
	Capture: func(HttpRequestEntity) {
	},
}

var _ = Describe("Token", func() {
	ctx := config.Context{}
	logger := lager.NewLogger("uaausersimport")
	ctx.Logger = logger
	ctx.TokenFetcher = &FakeTokenFetcher{"test-access-token"}
	Describe("Retrieve token from UAA", func() {
		NewGateway = func() HttpGateway {
			return mockGateway
		}
		Context("With correct response", func() {
			It("should get a valid token", func() {
				token, err := GetToken(&ctx)
				Ω(err).Should(BeNil())
				Ω(token).Should(Equal("test-access-token"))
			})
		})
	})
})
