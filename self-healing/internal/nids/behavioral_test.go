package nids

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func TestBehavioralEngine_PortScanDetection(t *testing.T) {
	be := NewBehavioralEngine()
	be.PortScanThresh = 5 // Low threshold for testing

	srcIP := net.ParseIP("192.168.1.50")
	dstIP := net.ParseIP("192.168.1.1")

	// Simulate 10 packets to different ports
	for i := 1; i <= 10; i++ {
		ip := &layers.IPv4{
			Version:  4,
			IHL:      5,
			TTL:      64,
			Protocol: layers.IPProtocolTCP,
			SrcIP:    srcIP.To4(),
			DstIP:    dstIP.To4(),
		}
		tcp := &layers.TCP{
			DstPort: layers.TCPPort(80 + i),
		}
		
		buffer := gopacket.NewSerializeBuffer()
		gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{}, ip, tcp, gopacket.Payload([]byte("probe")))
		packet := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
		
		be.AnalyzePacket(packet)
	}

	be.mu.Lock()
	stats := be.Stats[srcIP.String()]
	be.mu.Unlock()

	if len(stats.UniquePorts) < be.PortScanThresh {
		t.Errorf("Expected port scan detection, got only %d ports tracked", len(stats.UniquePorts))
	}
}
