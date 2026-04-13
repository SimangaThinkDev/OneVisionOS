package nids

import (
	"net"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestNIDSEngine_Detection(t *testing.T) {
	engine := NewEngine()
	engine.LoadSignatures(DefaultSignatures())

	// Create a mock IPv4 + TCP packet with malicious payload
	payload := []byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n/bin/bash -i")
	
	ip := &layers.IPv4{
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		SrcIP:    net.ParseIP("192.168.1.100").To4(),
		DstIP:    net.ParseIP("192.168.1.1").To4(),
	}
	tcp := &layers.TCP{
		SrcPort: 12345,
		DstPort: 80,
	}
	
	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	gopacket.SerializeLayers(buffer, opts,
		ip,
		tcp,
		gopacket.Payload(payload),
	)
	
	packet := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)

	// Process the packet
	go engine.ProcessPacket(packet)

	// Wait for alert
	select {
	case alert := <-engine.Alerts:
		if alert.Signature.ID != "OV-001" {
			t.Errorf("Expected signature OV-001, got %s", alert.Signature.ID)
		}
		t.Logf("Detected alert from source: %s", alert.SourceIP)
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout: No alert generated for malicious packet")
	}
}
