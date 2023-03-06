package irremote

import (
	"github.com/xkamail/huberlink-platform/pkg/snowid"
)

type MQTTCmd struct {
	Code      string `json:"code"`
	Frequency int    `json:"frequency"`
}

type MQTTReport struct {
	VirtualKeyID snowid.ID  `json:"virtualKeyId"`
	Properties   Properties `json:"properties"`
}
