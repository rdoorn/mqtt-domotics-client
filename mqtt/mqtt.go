package mqtt

import (
	"log"

	mqttc "github.com/eclipse/paho.mqtt.golang"
)

// Handler handles the MQTT connections
type Handler struct {
	Subscriptions chan mqttc.Message
	client        mqttc.Client
}

// New starts a new MQTT connection
func New(host string) (*Handler, error) {
	h := &Handler{
		Subscriptions: make(chan mqttc.Message),
	}
	if err := h.connect(host); err != nil {
		return nil, err
	}
	return h, nil
}

// Close closes a MQTT connection and its handlers
func (h *Handler) Close() {
	h.client.Disconnect(0)
}

// Subscribe adds a subscription
func (h *Handler) Subscribe(topic string) {
	h.subscribe(topic)
}

// Publish publishes data to MQTT
func (h *Handler) Publish(topic, message string) {
	go h.publish(topic, message)
}

func (h *Handler) subscribe(topic string) {
	if token := h.client.Subscribe(topic, 0, func(client mqttc.Client, msg mqttc.Message) {
		h.Subscriptions <- msg
		// wg.Done()
	}); token.Wait() && token.Error() != nil {
		log.Printf("subscriber error: %s", token.Error())
		return
	}
}

func (h *Handler) publish(topic, message string) {
	if token := h.client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		log.Printf("publish error: %s", token.Error())
		return
	}
}

func (h *Handler) connect(host string) error {
	opts := mqttc.NewClientOptions().AddBroker(host)

	h.client = mqttc.NewClient(opts)
	if token := h.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
