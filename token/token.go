package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-golang/lager"
	ghttp "github.com/pivotalservices/gtils/http"
	"github.com/pivotalservices/uaausersimport/config"
	"github.com/pivotalservices/uaausersimport/functions"
)

type Token struct {
	AccessToken string `json:"access_token"`
}

var NewGateway = func() ghttp.HttpGateway {
	return ghttp.NewHttpGateway()
}

var GetTokenOld functions.GetTokenFunc = func(c *config.Context) (token string, err error) {
	entity := ghttp.HttpRequestEntity{
		Url:      fmt.Sprintf("%s/oauth/token?grant_type=client_credentials", c.UAAURL),
		Username: c.Clientid,
		Password: c.Secret,
	}
	httpGateway := NewGateway()
	request := httpGateway.Post(entity, nil)
	response, err := request()
	if err != nil {
		return
	}
	return parse(response)
}

var GetToken functions.GetTokenFunc = func(c *config.Context) (tokenString string, err error) {
	c.Logger.Debug("GetTokenFunc", lager.Data{"msg": "Start fetch token............."})
	token, err := c.TokenFetcher.FetchToken(true)
	if err != nil {
		return
	}
	tokenString = token.AccessToken
	c.Logger.Debug("GetTokenFunc", lager.Data{"msg": "Finish fetch token", "token": tokenString})
	return
}

func parse(response *http.Response) (tokenString string, err error) {
	body := response.Body
	defer body.Close()
	token := &Token{}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return
	}
	tokenString = token.AccessToken
	return
}
