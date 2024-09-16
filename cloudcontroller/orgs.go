package cloudcontroller

import (
	"fmt"

	"github.com/pivotalservices/uaausersimport/config"
	"github.com/pivotalservices/uaausersimport/uaa"
)

type OrgInfo struct {
	uaa.UserIdInfo
	Org  config.Org
	Guid string
}

type AssociateOrgFunc func(OrgInfo) (string, error)

var AssociateOrg AssociateOrgFunc = func(info OrgInfo) (guid string, err error) {
	info.Logger.Debug("Invoking AssociateOrg on CloudController")
	response, err := info.RequestFn(info.Token, fmt.Sprintf("%s/v2/organizations?q=name:%s", info.Ccurl, info.Org.Name), "GET", "application/json", nil)
	if err != nil {
		return
	}
	guid, err = parseResponse(response)
	if err != nil {
		return
	}
	info.Logger.Debug(fmt.Sprintf("Associate user id :%s to org: %s.........", info.User.Uid, guid))
	_, err = info.RequestFn(info.Token, fmt.Sprintf("%s/v2/organizations/%s/users/%s", info.Ccurl, guid, info.UserId), "PUT", "application/json", nil)
	if err != nil {
		return
	}
	for _, role := range info.Org.Roles {
		info.Logger.Debug(fmt.Sprintf("Associate user id :%s to org: %s with %s role.........", info.User.Uid, guid, role))
		_, err = info.RequestFn(info.Token, fmt.Sprintf("%s/v2/organizations/%s/%s/%s", info.Ccurl, guid, role, info.UserId), "PUT", "application/json", nil)
		if err != nil {
			return
		}
	}
	info.Logger.Debug("Finish invoking AssociateOrg on CloudController")
	return
}
