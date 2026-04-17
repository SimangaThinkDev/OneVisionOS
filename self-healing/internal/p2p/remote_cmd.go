package p2p

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const CommandProtocolID = "/onevision/command/1.0.0"

type CommandRequest struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type CommandResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error,omitempty"`
}

// RegisterCommandProtocol registers the handler for remote commands
func (n *Node) RegisterCommandProtocol() {
	n.Host.SetStreamHandler(CommandProtocolID, n.handleCommandStream)
	log.Printf("[P2P] Remote command protocol registered: %s", CommandProtocolID)
}

// handleCommandStream handles incoming command execution requests
func (n *Node) handleCommandStream(s network.Stream) {
	defer s.Close()

	var req CommandRequest
	if err := json.NewDecoder(s).Decode(&req); err != nil {
		log.Printf("[P2P] Failed to decode command request: %v", err)
		return
	}

	log.Printf("[P2P] Executing remote command from %s: %s %v", s.Conn().RemotePeer(), req.Command, req.Args)

	// Execute the command
	cmd := exec.Command(req.Command, req.Args...)
	stdout, err := cmd.Output()
	
	resp := CommandResponse{
		Stdout: string(stdout),
	}

	if err != nil {
		resp.Error = err.Error()
		if exitErr, ok := err.(*exec.ExitError); ok {
			resp.Stderr = string(exitErr.Stderr)
			resp.ExitCode = exitErr.ExitCode()
		} else {
			resp.ExitCode = -1
		}
	} else {
		resp.ExitCode = 0
	}

	if err := json.NewEncoder(s).Encode(resp); err != nil {
		log.Printf("[P2P] Failed to send command response: %v", err)
	}
}

// SendCommand sends a command to a specific peer and returns the response
func (n *Node) SendCommand(peerID peer.ID, req CommandRequest) (*CommandResponse, error) {
	s, err := n.Host.NewStream(n.ctx, peerID, CommandProtocolID)
	if err != nil {
		return nil, fmt.Errorf("open stream: %w", err)
	}
	defer s.Close()

	if err := json.NewEncoder(s).Encode(req); err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	var resp CommandResponse
	if err := json.NewDecoder(s).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &resp, nil
}
