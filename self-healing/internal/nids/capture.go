// +build pcap

package nids

import (
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func (e *Engine) StartCapture(device string) {
	log.Printf("[NIDS] Starting live capture on %s", device)
	
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Printf("[NIDS] Error opening device %s: %v. Capture disabled.", device, err)
		return
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		e.ProcessPacket(packet)
	}
}

// StartAutoCapture attempts to find a default device and start capturing
func (e *Engine) StartAutoCapture() {
	devices, err := pcap.FindAllDevs()
	if err != nil || len(devices) == 0 {
		log.Printf("[NIDS] No network devices found for capture.")
		return
	}

	// Prefer eth0 or wlan0 or similar, otherwise just the first one
	target := devices[0].Name
	for _, d := range devices {
		if d.Name == "eth0" || d.Name == "en0" || d.Name == "wlan0" {
			target = d.Name
			break
		}
	}

	go e.StartCapture(target)
}
