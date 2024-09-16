package cloudcontroller

import (
	"fmt"

	"github.com/pivotalservices/uaausersimport/config"
)

type SpaceInfo struct {
	OrgInfo
	Space config.Space
}

type AssociateSpaceFunc func(SpaceInfo) error

var AssociateSpace AssociateSpaceFunc = func(info SpaceInfo) (err error) {
	info.Logger.Debug("Invoking AssociateSpace on CloudController")
	response, err := info.RequestFn(info.Token, fmt.Sprintf("%s/v2/spaces?q=name:%s&q=organization_guid:%s", info.Ccurl, info.Space.Name, info.OrgInfo.Guid), "GET", "application/json", nil)
	if err != nil {
		return
	}
	guid, err := parseResponse(response)
	if err != nil {
		return
	}
	for _, role := range info.Space.Roles {
		info.Logger.Debug(fmt.Sprintf("Associate user id :%s to space: %s with %s role.........", info.User.Uid, guid, role))
		_, err = info.RequestFn(info.Token, fmt.Sprintf("%s/v2/spaces/%s/%s/%s", info.Ccurl, guid, role, info.UserId), "PUT", "application/json", nil)
		if err != nil {
			return
		}
	}
	info.Logger.Debug("Finish invoking AssociateSpace on CloudController")
	return
}
