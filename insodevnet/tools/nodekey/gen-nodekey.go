package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

func main() {

	keyDir := filepath.Join("insodevnet/docker/single-validator/nodekey")

	if err := os.MkdirAll(keyDir, 0700); err != nil {
		log.Fatalf("❌ Failed to create key directory: %v", err)
	}

	privKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate nodekey: %v", err)
	}

	keyFile := filepath.Join(keyDir, "nodekey")

	// Save raw hex key
	raw := crypto.FromECDSA(privKey)
	nodeID := crypto.PubkeyToAddress(privKey.PublicKey).Hex()[2:]
	if err := os.WriteFile(keyFile, []byte(hex.EncodeToString(raw)), 0600); err != nil {
		log.Fatalf("Failed to write nodekey to %s: %v", keyFile, err)
	}

	pubkey := crypto.FromECDSAPub(&privKey.PublicKey)
	hexPub := hex.EncodeToString(pubkey)

	// Enode URL
	enodeURL := makeEnode(privKey, "127.0.0.1", 30303)

	fmt.Printf("✅ Nodekey saved: %s\n", keyFile)
	fmt.Printf("🔗 Enode URL: %s\n", enodeURL)
	fmt.Printf("🧪 Node ID: %s\n", nodeID)
	fmt.Printf("📎 Public key: %s\n", hexPub)

	// Save static-nodes.json
	staticNodes := []string{enodeURL}
	staticNodesPath := filepath.Join(keyDir, "static-nodes.json")
	staticData, _ := json.MarshalIndent(staticNodes, "", "  ")
	os.WriteFile(staticNodesPath, staticData, 0644)

	fmt.Printf("📄 static-nodes.json written: %s\n", staticNodesPath)
}

func makeEnode(priv *ecdsa.PrivateKey, ip string, port int) string {
	node := enode.NewV4(&priv.PublicKey, net.ParseIP(ip), int(uint16(port)), int(uint16(port)))
	return node.String()
}
