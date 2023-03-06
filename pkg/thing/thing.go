package thing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/slog"

	"github.com/xkamail/huberlink-platform/iot/device"
	"github.com/xkamail/huberlink-platform/iot/irremote"
	"github.com/xkamail/huberlink-platform/pkg/snowid"
)

func New() (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://mc.xkamail.me:1883")
	opts.SetPassword("test")
	opts.SetUsername("test")
	opts.SetClientID("")
	opts.SetCleanSession(true)
	opts.SetKeepAlive(60 * time.Second)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}

type ExecuteMessage struct {
	Time int64 `json:"time"`
	Data any   `json:"data"`
}

func Execute(ctx context.Context, deviceId snowid.ID, m *ExecuteMessage) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		c, err := New()
		if err != nil {
			return err
		}
		defer c.Disconnect(250)
		topic := fmt.Sprintf("huberlink/%s/thing/property/report", deviceId.String())

		c.Publish(topic, 1, false, string(b)).
			WaitTimeout(2 * time.Second)
		return nil
	}
}

type ReportFormat[T any] struct {
	Time int64 `json:"time"`
	Data T     `json:"data"`
}

func SubReport(ctx context.Context) mqtt.MessageHandler {
	fn := func(client mqtt.Client, msg mqtt.Message) error {
		topic := msg.Topic()
		deviceId, err := deviceIDFromTopic(topic)
		if err != nil {
			return err
		}
		fmt.Printf("device %s ", deviceId.String())
		d, err := device.Find(ctx, deviceId)
		if err != nil {
			return err
		}
		switch d.Kind {
		case device.KindIRRemote:
			var d ReportFormat[irremote.MQTTReport]
			if err := json.Unmarshal(msg.Payload(), &d); err != nil {
				return err
			}

			msg.Ack()

		default:
			return fmt.Errorf("invalid device kind %d", d.Kind)
		}
		return nil
	}
	return func(client mqtt.Client, msg mqtt.Message) {
		if err := fn(client, msg); err != nil {
			slog.Error("report error", err)
		}
	}
}

func deviceIDFromTopic(topic string) (snowid.ID, error) {
	if !strings.HasPrefix(topic, "huberlink/") {
		return snowid.Zero, fmt.Errorf("invalid topic %s", topic)
	}
	s := strings.Split(topic, "/")
	if len(s) < 2 {
		return snowid.Zero, fmt.Errorf("invalid topic %s", topic)
	}
	return snowid.Parse(s[1])
}
