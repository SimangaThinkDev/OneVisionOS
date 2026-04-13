package nids

import (
	"log"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TrafficStats struct {
	PacketCount int
	ByteCount   int
	UniquePorts map[int]struct{}
	LastSeen    time.Time
}

type BehavioralEngine struct {
	Stats        map[string]*TrafficStats
	mu           sync.Mutex
	Threshold    int // packets per second per IP
	PortScanThresh int // unique ports per window
}

func NewBehavioralEngine() *BehavioralEngine {
	return &BehavioralEngine{
		Stats:          make(map[string]*TrafficStats),
		Threshold:      500, // example threshold
		PortScanThresh: 20,
	}
}

func (b *BehavioralEngine) AnalyzePacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)
	srcIP := ip.SrcIP.String()

	b.mu.Lock()
	defer b.mu.Unlock()

	stats, exists := b.Stats[srcIP]
	if !exists {
		stats = &TrafficStats{
			UniquePorts: make(map[int]struct{}),
		}
		b.Stats[srcIP] = stats
	}

	stats.PacketCount++
	stats.ByteCount += len(packet.Data())
	stats.LastSeen = time.Now()

	// Track unique destination ports for port scan detection
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		stats.UniquePorts[int(tcp.DstPort)] = struct{}{}
	}

	// Simple Threshold Check
	if stats.PacketCount > b.Threshold {
		log.Printf("[SECURITY] ANOMALY: High traffic volume from %s (%d packets)", srcIP, stats.PacketCount)
		// Trigger Isolation (Phase 14-Task 4)
	}

	if len(stats.UniquePorts) > b.PortScanThresh {
		log.Printf("[SECURITY] ANOMALY: Potential port scanning from %s (%d ports)", srcIP, len(stats.UniquePorts))
	}
}

func (b *BehavioralEngine) ResetStats() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			b.mu.Lock()
			// Reset or decay stats
			for ip, s := range b.Stats {
				if time.Since(s.LastSeen) > 30*time.Second {
					delete(b.Stats, ip)
				} else {
					s.PacketCount = 0
					s.UniquePorts = make(map[int]struct{})
				}
			}
			b.mu.Unlock()
		}
	}()
}
