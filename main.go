package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	uaa "github.com/cloudfoundry-incubator/uaa-token-fetcher"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/pivotalservices/uaausersimport/config"
	"github.com/pivotalservices/uaausersimport/token"
	"gopkg.in/yaml.v2"
)

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func setup() *config.Context {
	logger := lager.NewLogger("uaausersimport")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	ctx := &config.Context{}
	ctx.Logger = logger

	err := parseEnv(ctx)
	fatalIf(err)

	return ctx
}

func main() {
	ctx := setup()
	run(ctx)
}

func run(ctx *config.Context) {
	err := token.GetToken.MapUsers().AddUAAUsers().AddCCUsers().MapOrgs().MapSpaces(ctx)
	fatalIf(err)
}

func parseEnv(ctx *config.Context) error {
	users := os.Getenv("USERS_CONFIG_FILE")
	file, err := os.Open(users)
	if err != nil {
		err = fmt.Errorf("Cannot open %s : %s", users, err.Error())
		return err
	}

	cfg, err := config.Parse(file)
	fatalIf(err)

	env := os.Getenv("CF_ENVIRONMENT")
	file, err = os.Open(env)
	if err != nil {
		err = fmt.Errorf("Cannot open %s : %s", env, err.Error())
		return err
	}

	data, err := ioutil.ReadAll(file)
	fatalIf(err)

	err = yaml.Unmarshal(data, ctx)
	fatalIf(err)

	if len(ctx.UAAURL) == 0 {
		ctx.UAAURL = os.Getenv("UAA_ADDRESS")
	}

	if len(ctx.Clientid) == 0 {
		ctx.Clientid = os.Getenv("CLIENT_ID")
	}

	if len(ctx.Secret) == 0 {
		ctx.Secret = os.Getenv("CLIENT_SECRET")
	}

	ctx.RequestFn = token.RequestWithToken
	ctx.Users = cfg.Users
	ctx.Origin = cfg.Origin

	tokenFetcherConfig := uaa.TokenFetcherConfig{
		MaxNumberOfRetries:                3,
		RetryInterval:                     15 * time.Second,
		ExpirationBufferTime:              30,
		DisableTLSCertificateVerification: true,
	}
	oauth := &uaa.OAuthConfig{
		TokenEndpoint: ctx.UAAURL,
		ClientName:    ctx.Clientid,
		ClientSecret:  ctx.Secret,
		Port:          443,
	}
	clk := clock.NewClock()
	fetcher, err := uaa.NewTokenFetcher(ctx.Logger, oauth, tokenFetcherConfig, clk)
	fatalIf(err)
	ctx.TokenFetcher = fetcher
	return nil
}
