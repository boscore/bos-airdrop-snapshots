package pkg

import (
	"log"
	"testing"
)

func TestActionsCache(t *testing.T) {
	actions, _ := GetCachedActions("updateauth")
	log.Printf("todo: %d", len(actions))
}
