package home

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/xkamail/huberlink-platform/pkg/snowid"
	"github.com/xkamail/huberlink-platform/pkg/uierr"
)

type SceneExecuter interface {
	ExecuteScene(scene Scene) error
}
type SceneRun uint

const (
	SceneRunWhenTap SceneRun = 1 << iota
	SceneRunWhenSchedule
)

type SceneScheduleRepeat uint

const (
	SceneScheduleRepeatSun = 1 << iota
	SceneScheduleRepeatMon
	SceneScheduleRepeatTue
	SceneScheduleRepeatWed
	SceneScheduleRepeatThu
	SceneScheduleRepeatFri
	SceneScheduleRepeatSat
)

type SceneScheduleTime string

func (t *SceneScheduleTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(*t))
}

func (t *SceneScheduleTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	// 00:00 is valid
	if s == "00:00" {
		*t = SceneScheduleTime(s)
		return nil
	}
	// hh:mm
	if len(s) != 5 {
		return uierr.BadInput("time", "invalid time format")
	}
	if s[2] != ':' {
		return uierr.BadInput("time", "invalid time format")
	}
	hh, err := strconv.ParseInt(s[:2], 10, 64)
	if err != nil {
		return uierr.BadInput("time", "invalid time format: hh should be 0-23")
	}
	// validate hh should be 0-23
	mm, err := strconv.ParseInt(s[3:], 10, 64)
	if err != nil {
		return uierr.BadInput("time", "invalid time format: mm should be 0-59")
	}
	// validate mm should be 0-59
	if hh < 0 || hh > 23 || mm < 0 || mm > 59 {
		return uierr.BadInput("time", "invalid time format: hh:mm should be 0-23:0-59")
	}

	*t = SceneScheduleTime(s)
	return nil
}

var _ json.Marshaler = (*SceneScheduleTime)(nil)
var _ json.Unmarshaler = (*SceneScheduleTime)(nil)

type Scene struct {
	ID     snowid.ID `json:"id"`
	HomeID snowid.ID `json:"homeId"`
	Name   string    `json:"name"`
	Run    SceneRun  `json:"run"`
	// when no repeat is zero value
	// multiple days can be set by bitwise or
	ScheduleRepeat SceneScheduleRepeat `json:"schedule"`
	// hh:mm format
	ScheduleTime SceneScheduleTime `json:"scheduleTime"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SceneAction struct {
	ID       snowid.ID `json:"id"`
	Order    uint      `json:"order"`
	Name     string    `json:"name"`
	DeviceID snowid.ID `json:"deviceId"` // device id to execute
	// action to execute
	RawAction string `json:"rawAction"`
}
