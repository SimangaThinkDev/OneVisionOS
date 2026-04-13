package nids

import (
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Engine struct {
	Signatures []Signature
	Alerts     chan Alert
}

func NewEngine() *Engine {
	return &Engine{
		Signatures: []Signature{},
		Alerts:     make(chan Alert, 100),
	}
}

func (e *Engine) LoadSignatures(sigs []Signature) {
	e.Signatures = sigs
}

func (e *Engine) ProcessPacket(packet gopacket.Packet) {
	// Extract IPv4 layer
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer == nil {
		return
	}
	ip, _ := ipLayer.(*layers.IPv4)

	// Extract TCP layer
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	var destPort int
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		destPort = int(tcp.DstPort)
	}

	// Extract payload
	var payload string
	appLayer := packet.ApplicationLayer()
	if appLayer != nil {
		payload = string(appLayer.Payload())
	} else if len(packet.Data()) > 0 {
		// Fallback to searching the entire packet data if application layer isn't decoded
		payload = string(packet.Data())
	} else {
		return
	}

	// Check against signatures
	for _, sig := range e.Signatures {
		match := false

		// Check port if specified
		if sig.DestPort != 0 && sig.DestPort != destPort {
			continue
		}

		// Check payload pattern
		if sig.PayloadPattern != "" && strings.Contains(payload, sig.PayloadPattern) {
			match = true
		}

		if match {
			alert := Alert{
				Signature: sig,
				SourceIP:  ip.SrcIP.String(),
				DestIP:    ip.DstIP.String(),
				Payload:   payload,
			}
			e.Alerts <- alert
			log.Printf("[NIDS] ALERT [%s]: %s detected from %s", sig.Level, sig.Name, ip.SrcIP)
		}
	}
}
