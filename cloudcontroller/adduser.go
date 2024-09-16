package cloudcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pivotalservices/uaausersimport/uaa"
)

type UserRequest struct {
	Id string `json:"guid"`
}

type AddCCUserFunc func(uaa.UserIdInfo) error

var AddCCUser AddCCUserFunc = func(info uaa.UserIdInfo) (err error) {
	info.Logger.Debug("Invoking AddCCUserFunc on CloudController")
	info.Logger.Debug(fmt.Sprintf("add user id: %s with Uaa id: %s .........", info.User.Uid, info.UserId))
	userRequest := UserRequest{
		Id: info.UserId,
	}
	data, err := json.Marshal(userRequest)
	if err != nil {
		return
	}
	body := bytes.NewBuffer(data)
	_, err = info.RequestFn(info.Token, fmt.Sprintf("%s/v2/users", info.Ccurl), "POST", "application/json", body)
	info.Logger.Debug("Finish invoking AddCCUserFunc on CloudController")
	return err
}
