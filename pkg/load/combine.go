package load

func BuildInternals() error {
	err := TriggerBuild("./sitedata/theme/", "./static/")
	if err != nil {
		return err
	}

	err = BuildConfigCache()
	if err != nil {
		return err
	}

	err = BuildResourceCache()
	if err != nil {
		return err
	}

	return nil
}
