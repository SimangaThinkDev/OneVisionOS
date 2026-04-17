package p2p

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const FileDistProtocolID = "/onevision/filedist/1.0.0"

type FileDistRequest struct {
	DestPath string `json:"dest_path"`
	Content  []byte `json:"content"`
	Mode     uint32 `json:"mode"`
}

type FileDistResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// RegisterFileDistProtocol registers the handler for incoming file distributions
func (n *Node) RegisterFileDistProtocol() {
	n.Host.SetStreamHandler(FileDistProtocolID, n.handleFileDistStream)
	log.Printf("[P2P] File distribution protocol registered: %s", FileDistProtocolID)
}

func (n *Node) handleFileDistStream(s network.Stream) {
	defer s.Close()

	var req FileDistRequest
	if err := json.NewDecoder(s).Decode(&req); err != nil {
		log.Printf("[P2P] Failed to decode file dist request: %v", err)
		return
	}

	resp := FileDistResponse{Success: true}

	if err := os.MkdirAll(filepath.Dir(req.DestPath), 0755); err != nil {
		resp.Success = false
		resp.Error = err.Error()
	} else if err := os.WriteFile(req.DestPath, req.Content, os.FileMode(req.Mode)); err != nil {
		resp.Success = false
		resp.Error = err.Error()
	} else {
		log.Printf("[P2P] File distributed to %s from %s", req.DestPath, s.Conn().RemotePeer())
	}

	json.NewEncoder(s).Encode(resp)
}

// DistributeFile sends a file to a specific peer
func (n *Node) DistributeFile(peerID peer.ID, destPath string, content []byte, mode uint32) error {
	s, err := n.Host.NewStream(n.ctx, peerID, FileDistProtocolID)
	if err != nil {
		return fmt.Errorf("open stream: %w", err)
	}
	defer s.Close()

	req := FileDistRequest{DestPath: destPath, Content: content, Mode: mode}
	if err := json.NewEncoder(s).Encode(req); err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	var resp FileDistResponse
	if err := json.NewDecoder(s).Decode(&resp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("peer error: %s", resp.Error)
	}
	return nil
}
