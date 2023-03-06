package irremote

type Air struct {
}

func (a Air) PropertiesKeys() []string {
	return []string{"power", "mode", "temp", "fan", "swing"}
}

func (a Air) Properties() Properties {
	return Properties{
		"power": "on",
		"mode":  "cool",
		"temp":  "25",
		"fan":   "auto",
		"swing": "off",
	}
}

var _ Virtualer = (*Air)(nil)
