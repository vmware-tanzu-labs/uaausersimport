package functions

import (
	cc "github.com/pivotalservices/uaausersimport/cloudcontroller"
	"github.com/pivotalservices/uaausersimport/config"
	uaa "github.com/pivotalservices/uaausersimport/uaa"
)

// GetTokenFunc DOCUMENT ME!
type GetTokenFunc func(*config.Context) (string, error)

// MapUsersFunc DOCUMENT ME!
type MapUsersFunc func(*config.Context) ([]uaa.UserInfo, error)

// AddUserFunc DOCUMENT ME!
type AddUserFunc func(*config.Context) ([]uaa.UserIdInfo, error)

// MapOrgsFunc DOCUMENT ME!
type MapOrgsFunc func(*config.Context) ([]cc.OrgInfo, error)

// MapUsers DOCUMENT ME!
func (getToken GetTokenFunc) MapUsers() MapUsersFunc {
	return func(ctx *config.Context) ([]uaa.UserInfo, error) {
		ctx.Logger.Debug("Start MapUsers")
		token, err := getToken(ctx)
		if err != nil {
			return nil, err
		}
		var users []uaa.UserInfo
		for _, user := range ctx.Users {
			userInfo := uaa.UserInfo{
				Context: ctx,
				Token:   token,
				User:    user,
				Origin:  ctx.Origin,
			}
			users = append(users, userInfo)
		}
		ctx.Logger.Debug("Finish MapUsers")
		return users, nil
	}
}

// AddUAAUsers DOCUMENT ME!
func (mapUsers MapUsersFunc) AddUAAUsers() AddUserFunc {
	return func(ctx *config.Context) ([]uaa.UserIdInfo, error) {
		ctx.Logger.Debug("Start AddUAAUsers")
		userInfos, err := mapUsers(ctx)
		if err != nil {
			return nil, err
		}
		var users []uaa.UserIdInfo
		for _, user := range userInfos {
			id, err := uaa.AddUAAUser(user)
			if err != nil {
				return nil, err
			}
			userIDInfo := uaa.UserIdInfo{
				UserInfo: user,
				UserId:   id,
			}
			users = append(users, userIDInfo)
		}
		ctx.Logger.Debug("Finish AddUAAUsers")
		return users, nil
	}
}

// AddCCUsers DOCUMENT ME!
func (addUsers AddUserFunc) AddCCUsers() AddUserFunc {
	return func(ctx *config.Context) ([]uaa.UserIdInfo, error) {
		ctx.Logger.Debug("Start AddCCUser")
		userIDInfos, err := addUsers(ctx)
		if err != nil {
			return nil, err
		}
		for _, user := range userIDInfos {
			err := cc.AddCCUser(user)
			if err != nil {
				return nil, err
			}
		}
		ctx.Logger.Debug("Finish AddCCUser")
		return userIDInfos, nil
	}
}

// MapOrgs DOCUMENT ME!
func (addUsers AddUserFunc) MapOrgs() MapOrgsFunc {
	return func(ctx *config.Context) ([]cc.OrgInfo, error) {
		ctx.Logger.Debug("Start MapOrgs")
		userIDInfos, err := addUsers(ctx)
		if err != nil {
			return nil, err
		}
		var orgInfos []cc.OrgInfo
		for _, userIDInfo := range userIDInfos {
			for _, org := range userIDInfo.User.Orgs {
				orgInfo := cc.OrgInfo{
					UserIdInfo: userIDInfo,
					Org:        org,
				}
				guid, err := cc.AssociateOrg(orgInfo)
				if err != nil {
					return nil, err
				}
				orgInfo.Guid = guid
				orgInfos = append(orgInfos, orgInfo)
			}
		}
		ctx.Logger.Debug("Finish MapOrgs")
		return orgInfos, nil
	}
}

// MapSpaces DOCUMENT ME!
func (mapOrgs MapOrgsFunc) MapSpaces(ctx *config.Context) error {
	return func(ctx *config.Context) error {
		ctx.Logger.Debug("Start MapSpaces")
		orgInfos, err := mapOrgs(ctx)
		if err != nil {
			return err
		}
		for _, orgInfo := range orgInfos {
			spaces := orgInfo.Org.Spaces
			for _, space := range spaces {
				spaceInfo := cc.SpaceInfo{
					Space:   space,
					OrgInfo: orgInfo,
				}
				err := cc.AssociateSpace(spaceInfo)
				if err != nil {
					return err
				}
			}
		}
		ctx.Logger.Debug("Finish MapSpaces")
		return nil
	}(ctx)
}
