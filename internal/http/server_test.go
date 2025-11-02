package httpserver

import (
	"testing"
)

func TestGetLocalIPs(t *testing.T) {
	ips := getLocalIPs()

	// We should get at least 0 IPs (could be 0 in some environments)
	// This is a basic sanity check
	if ips == nil {
		t.Error("getLocalIPs returned nil, expected empty slice or IPs")
	}

	// If we got IPs, verify they are not empty and not loopback
	for _, ip := range ips {
		if ip == "" {
			t.Error("Got empty IP address")
		}
		if ip == "127.0.0.1" {
			t.Error("Got loopback IP, should have been filtered out")
		}
	}
}
