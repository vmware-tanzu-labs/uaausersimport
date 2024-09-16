package config

import (
	"io"
	"io/ioutil"
	"net/http"

	uaa "github.com/cloudfoundry-incubator/uaa-token-fetcher"
	"github.com/pivotal-golang/lager"
	"gopkg.in/yaml.v2"
)

type Space struct {
	Name  string
	Roles []string
}

type Org struct {
	Name   string
	Roles  []string
	Spaces []Space
}

type User struct {
	Uid        string
	Externalid string
	Emails     []string
	Orgs       []Org
}

type Config struct {
	Origin string
	Users  []User
}

// Context DOCUMENT ME!
type Context struct {
	Ccurl        string
	UAAURL       string
	Clientid     string
	Secret       string
	Origin       string
	RequestFn    func(string, string, string, string, io.Reader) (*http.Response, error)
	TokenFetcher uaa.TokenFetcher
	Users        []User
	Logger       lager.Logger
}

func Parse(reader io.Reader) (*Config, error) {
	config := Config{}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}
