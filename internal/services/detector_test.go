package services

import (
	"testing"
)

func TestDetectorService(t *testing.T) {
	svc := NewDetectorService()
	clients := svc.DetectClients(false)
	t.Logf("Found %d clients:", len(clients))
	for i, c := range clients {
		t.Logf("[%d] %s (%s, Mods: %d)", i+1, c.Path, c.Source, c.ModCount)
	}
}
