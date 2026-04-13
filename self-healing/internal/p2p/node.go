package p2p

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const (
	DiscoveryServiceTag = "onevision-p2p"
	HealthProtocolID     = "/onevision/health/1.2.0"
)

type Node struct {
	Host host.Host
	ctx  context.Context
}

// NewNode initializes a new libp2p host
func NewNode(ctx context.Context) (*Node, error) {
	// Create a libp2p Host with default settings and a random listening port
	h, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
	)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Host: h,
		ctx:  ctx,
	}

	// Set a stream handler for the Health Protocol
	node.Host.SetStreamHandler(HealthProtocolID, node.handleHealthStream)

	log.Printf("[P2P] Node started with Peer ID: %s", h.ID().String())
	log.Printf("[P2P] Listening on: %v", h.Addrs())

	return node, nil
}

// SetupDiscovery sets up mDNS discovery
func (n *Node) SetupDiscovery() error {
	ser := mdns.NewMdnsService(n.Host, DiscoveryServiceTag, &discoveryNotifee{h: n.Host})
	return ser.Start()
}

// handleHealthStream handles incoming health status updates
func (n *Node) handleHealthStream(s network.Stream) {
	defer s.Close()

	var status HealthStatus
	if err := json.NewDecoder(s).Decode(&status); err != nil {
		log.Printf("[P2P] Failed to decode health status from %s: %v", s.Conn().RemotePeer(), err)
		return
	}

	log.Printf("[P2P] Received health status from %s: %s (at %s)", status.Hostname, status.Status, status.Timestamp.Format(time.Kitchen))
}

// BroadcastHealthStatus sends the current node's health status to all connected peers
func (n *Node) BroadcastHealthStatus(status string) {
	hostname, _ := os.Hostname()
	msg := HealthStatus{
		PeerID:    n.Host.ID().String(),
		Hostname:  hostname,
		Status:    status,
		Timestamp: time.Now(),
	}

	peers := n.Host.Network().Peers()
	if len(peers) == 0 {
		return
	}

	for _, p := range peers {
		go func(peerID peer.ID) {
			s, err := n.Host.NewStream(n.ctx, peerID, HealthProtocolID)
			if err != nil {
				log.Printf("[P2P] Failed to open stream to %s: %v", peerID, err)
				return
			}
			defer s.Close()

			if err := json.NewEncoder(s).Encode(msg); err != nil {
				log.Printf("[P2P] Failed to send health status to %s: %v", peerID, err)
			}
		}(p)
	}
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// Connect to the peer
	if pi.ID == n.h.ID() {
		return // Skip self
	}

	log.Printf("[P2P] Found peer: %s. Connecting...", pi.ID.String())
	if err := n.h.Connect(context.Background(), pi); err != nil {
		log.Printf("[P2P] Connection failed for %s: %v", pi.ID.String(), err)
	} else {
		log.Printf("[P2P] Successfully connected to %s", pi.ID.String())
	}
}
