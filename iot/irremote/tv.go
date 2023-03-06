package irremote

type TV struct {
}

// Properties default properties of TV
func (T TV) Properties() Properties {
	return Properties{
		"power": "on",
		"mute":  "off",
		"vol":   "25",
		"ch":    "1",
	}
}

func (T TV) PropertiesKeys() []string {
	return []string{"power", "mute", "vol", "ch"}
}

var _ Virtualer = (*TV)(nil)
