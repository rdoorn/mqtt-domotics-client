package mqtt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	topic   = "testtopic"
	message = "hello world"
)

func TestMQ(t *testing.T) {
	c, err := New("tcp://localhost:1883")
	assert.Nil(t, err)

	c.Subscribe(topic)

	c.Publish(topic, message)

	timeout := time.NewTimer(10 * time.Second).C

	for {
		select {
		case msg := <-c.Subscriptions:
			assert.Equal(t, message, string(msg.Payload()))
			return
		case <-timeout:
			assert.Error(t, fmt.Errorf("Timeout received on subscription"))
		}
	}

}
