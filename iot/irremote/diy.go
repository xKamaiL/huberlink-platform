package irremote

type Other struct {
}

func (d Other) Properties() Properties {
	return make(Properties)
}

func (d Other) PropertiesKeys() []string {
	// return empty mean no validation
	return make([]string, 0)
}
