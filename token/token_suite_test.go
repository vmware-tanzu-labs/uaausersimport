package token_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/uaa-token-fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Token Suite")
}

type FakeTokenFetcher struct {
	AccessToken string
}

func (fake *FakeTokenFetcher) FetchToken(useCachedToken bool) (*token_fetcher.Token, error) {
	return &token_fetcher.Token{
		AccessToken: fake.AccessToken,
	}, nil
}
