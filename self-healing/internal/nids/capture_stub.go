// +build !pcap

package nids

import (
	"log"
)

func (e *Engine) StartCapture(device string) {
	log.Printf("[NIDS] Live capture is disabled (compiled without libpcap)")
}

func (e *Engine) StartAutoCapture() {
	log.Printf("[NIDS] Auto-capture is disabled (compiled without libpcap)")
}
