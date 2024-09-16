package uaa

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pivotalservices/uaausersimport/config"
)

type Email struct {
	Value string `json:"value"`
}

type UserRequest struct {
	UserName   string  `json:"userName"`
	Emails     []Email `json:"emails"`
	Origin     string  `json:"origin"`
	Externalid string  `json:"externalid"`
}

type UserCreateResponse struct {
	Id string
}

type UserInfo struct {
	*config.Context
	Token  string
	User   config.User
	Origin string
}

type UserIdInfo struct {
	UserInfo
	UserId string
}

type AddUAAUserFunc func(UserInfo) (string, error)

var NewRoundTripper = func() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

var AddUAAUser AddUAAUserFunc = func(info UserInfo) (userId string, err error) {
	info.Logger.Debug("Invoking AddUAAUserFunc")
	info.Logger.Debug(fmt.Sprintf("add user id: %s .........", info.User.Uid))
	emails := make([]Email, 0)
	for _, value := range info.User.Emails {
		email := Email{
			Value: value,
		}
		emails = append(emails, email)
	}
	userRequest := UserRequest{
		UserName:   info.User.Uid,
		Externalid: info.User.Externalid,
		Emails:     emails,
		Origin:     info.Origin,
	}
	data, err := json.Marshal(userRequest)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(data)
	response, err := info.RequestFn(info.Token, fmt.Sprintf("%s/Users", info.UAAURL), "POST", "application/json", body)
	if err != nil {
		return
	}
	info.Logger.Debug("Finish invoking AddUAAUserFunc")
	return parse(response)
}

func parse(response *http.Response) (guid string, err error) {
	body := response.Body
	defer body.Close()
	uResponse := UserCreateResponse{}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &uResponse)
	if err != nil {
		return
	}
	guid = uResponse.Id
	return
}
