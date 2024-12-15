package load

import "github.com/zulubit/mimi/pkg/read"

type ResoruceMap *[]read.Resource

var config *read.Config
var resources ResoruceMap

func BuildConfigCache() error {

	rc, err := read.ReadConfig("./sitedata/config.json")
	if err != nil {
		return err
	}

	config = rc

	return nil
}

func GetConfig() (*read.Config, error) {

	if config == nil {
		err := BuildConfigCache()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func BuildResourceCache() error {

	rc, err := read.ReadResources("./sitedata/resources")
	if err != nil {
		return err
	}
	resources = rc

	return nil
}

func GetResources() (ResoruceMap, error) {
	if resources == nil {
		err := BuildResourceCache()
		if err != nil {
			return nil, err
		}
	}

	return resources, nil
}
