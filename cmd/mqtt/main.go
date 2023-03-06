package main

import (
	"context"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/xkamail/huberlink-platform/iot/irremote"
	"github.com/xkamail/huberlink-platform/pkg/snowid"
	"github.com/xkamail/huberlink-platform/pkg/thing"
)

var _ mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run MQTT server
func run() error {
	c, err := thing.New()
	if err != nil {
		return err
	}
	err = thing.Execute(context.Background(), snowid.Gen(), &thing.ExecuteMessage{
		Time: time.Now().Unix(),
		Data: &irremote.MQTTCmd{
			Code:      "0x123",
			Frequency: 38,
		},
	})
	if err != nil {
		return err
	}
	c.IsConnected()
	log.Printf("server is running %v\n", c.IsConnected())
	select {}
}

/*
 tylink/${deviceId}/thing/property/report
*/
