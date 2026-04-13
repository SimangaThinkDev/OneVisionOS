package p2p

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const FileSyncProtocolID = "/onevision/filesync/1.0.0"

// RegisterFileSyncProtocol registers the handler for file requests
func (n *Node) RegisterFileSyncProtocol() {
	n.Host.SetStreamHandler(FileSyncProtocolID, func(s network.Stream) {
		defer s.Close()

		// Read the requested path
		buf := make([]byte, 1024)
		nRead, err := s.Read(buf)
		if err != nil {
			log.Printf("[P2P] FileSync: Failed to read request: %v", err)
			return
		}
		path := string(buf[:nRead])

		log.Printf("[P2P] FileSync: Peer %s requested file %s", s.Conn().RemotePeer(), path)

		// Basic security check: in production we would verify against manifest
		file, err := os.Open(path)
		if err != nil {
			log.Printf("[P2P] FileSync: File %s not available: %v", path, err)
			return
		}
		defer file.Close()

		_, err = io.Copy(s, file)
		if err != nil {
			log.Printf("[P2P] FileSync: Failed to send file %s: %v", path, err)
		}
	})
}

// FetchFileFromPeer attempts to download a file from a specific peer
func (n *Node) FetchFileFromPeer(ctx context.Context, peerID peer.ID, path string, targetPath string) error {
	s, err := n.Host.NewStream(ctx, peerID, FileSyncProtocolID)
	if err != nil {
		return err
	}
	defer s.Close()

	// Send requested path
	_, err = s.Write([]byte(path))
	if err != nil {
		return err
	}

	// Create target file
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, s)
	return err
}

// FetchFileFromAnyPeer attempts to download a file from any available peer
func (n *Node) FetchFileFromAnyPeer(ctx context.Context, path string, targetPath string) error {
	peers := n.Host.Network().Peers()
	if len(peers) == 0 {
		return fmt.Errorf("no connected peers available")
	}

	for _, p := range peers {
		log.Printf("[P2P] Attempting to fetch %s from peer %s", path, p)
		err := n.FetchFileFromPeer(ctx, p, path, targetPath)
		if err == nil {
			return nil
		}
		log.Printf("[P2P] Failed to fetch from peer %s: %v", p, err)
	}

	return fmt.Errorf("file not found on any connected peers")
}
