package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var config serverConfig

type serverConfig struct {
	ServerAddr      string   `yaml:"server_addr"`
	ServerBannedIPs []string `yaml:"server_banned_ips"`
	Jwtsecret       string   `yaml:"jwtsecret"`
	ServerUrl       string   `yaml:"server_url"`

	ChallengesDir string `yaml:"challenges_dir"`
	SlidesDir     string `yaml:"slides_dir"`
	FrontendDir   string `yaml:"frontend_dir"`

	AdminPassword string `yaml:"admin_password"`
	UserPassword  string `yaml:"user_password"`
	GoogleId      string `yaml:"googleId"`
	GoogleSecret  string `yaml:"googleSecret"`

	ContainerHosts []sContainerHost `yaml:"container_hosts"`
	BaseContainers []sBaseContainer `yaml:"base_containers"`
}

type sBaseContainer struct {
	Id   string
	Name string
	Bits string
}

type sContainerHost struct {
	HostnameAlias string
	Hostname      string
	Aslr          bool
	Arch          string
	SshBasePort   int
}

func parseConfig() error {
	data, err := ioutil.ReadFile("yookiterm-server.yml")
	if os.IsNotExist(err) {
		return fmt.Errorf("the configuration file (yookiterm-server.yml) doesn't exist")
	} else if err != nil {
		return fmt.Errorf("unable to read the configuration %s", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("unable to parse the configuration %s", err)
	}

	if config.ServerAddr == "" {
		config.ServerAddr = ":8080"
	}

	return nil
}

func getBaseContainerByName(containerBaseName string) sBaseContainer {
	for _, element := range config.BaseContainers {
		if element.Name == containerBaseName {
			return element
		}
	}

	return sBaseContainer{}
}

func getContainerHostByAlias(hostnameAlias string) sContainerHost {
	for _, element := range config.ContainerHosts {
		if element.HostnameAlias == hostnameAlias {
			return element
		}
	}

	return sContainerHost{}
}
