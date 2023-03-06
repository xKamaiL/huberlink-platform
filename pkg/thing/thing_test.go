package thing

import (
	"fmt"
	"testing"

	"github.com/xkamail/huberlink-platform/pkg/snowid"
)

func Test_deviceIDFromTopic(t *testing.T) {
	type args struct {
		topic string
	}
	id := snowid.Gen()
	tests := []struct {
		name    string
		args    args
		want    snowid.ID
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic: fmt.Sprintf("huberlink/%s/thing/property/report", id.String()),
			},
			want:    id,
			wantErr: false,
		},
		{
			"error",
			args{
				topic: "huberlink/1x/thing/property/report/1",
			},
			snowid.Zero,
			true,
		},
		{
			"not start with huberlink",
			args{
				topic: "huberlinkx/1/thing/property/report/1",
			},
			snowid.Zero,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deviceIDFromTopic(tt.args.topic)
			if (err != nil) != tt.wantErr {
				t.Errorf("deviceIDFromTopic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("deviceIDFromTopic() got = %v, want %v", got, tt.want)
			}
		})
	}
}
