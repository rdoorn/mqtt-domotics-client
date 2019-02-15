package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdoorn/mqtt-domotics-client/mqtt"
)

const (
	RECEIVE = "domoticz/out"
	SEND    = "domoticz/in"
	HOST    = "tcp://localhost:1883"
)

func main() {

	// wait for sigint or sigterm for cleanup - note that sigterm cannot be caught
	sigterm := make(chan os.Signal, 10)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	c, err := mqtt.New(HOST)
	if err != nil {
		log.Fatal(err.Error())
	}

	c.Subscribe(RECEIVE)

	for {
		select {
		case <-sigterm:
			log.Print("Exiting...")
			c.Close()
			return
		case msg := <-c.Subscriptions:
			log.Printf("message received [%s]: %s", msg.Topic(), msg.Payload())
		}
	}
}

/*
func connect() mqtt.Client {
	opts := mqtt.NewClientOptions().AddBroker(HOST)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	return client
}

type Handler struct {
  Receive chan mqtt.Message
  Send chan mqtt.Message

}

func main() {

	client := connect()
  subscribe()


	var wg sync.WaitGroup
	wg.Add(1)

	log.Printf("Subscribed to: %s", TOPIC)

	if token := client.Subscribe(TOPIC, 0, func(client mqtt.Client, msg mqtt.Message) {

		log.Printf("got message: %+v", msg)

		   if string(msg.Payload()) != "mymessage" {
		           t.Fatalf("want mymessage, got %s", msg.Payload())
		   }

		wg.Done()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}


	  if token := client.Publish(TOPIC, 0, false, "mymessage"); token.Wait() && token.Error() != nil {
	          t.Fatal(token.Error())
	  }

	wg.Wait()
}
*/
