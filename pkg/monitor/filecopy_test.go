package monitor

import (
	"github.com/rjeczalik/notify"
	"log"
	"testing"
)

func TestStart(t *testing.T) {
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch("/home/tangxu/file-copy-test", c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Block until an event is received.
	ei := <-c
	log.Println("Got event:", ei)
}
