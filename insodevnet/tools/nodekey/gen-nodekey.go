package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"log"
	"net"
	"os"
)

func main() {
	dir := "insodevnet/keys/nodekeystore"

	// Create keys dir if missing
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			log.Fatalf("Failed to create keys dir: %v", err)
		}
	}

	// Generate private key
	privKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate nodekey: %v", err)
	}

	node := enode.NewV4(&privKey.PublicKey, nil, 0, 0)
	keyFileName := fmt.Sprintf("nodekey_%s", node.ID().String())
	keyFile := fmt.Sprintf("%s/%s", dir, keyFileName)

	raw := crypto.FromECDSA(privKey)
	if err := os.WriteFile(keyFile, []byte(hex.EncodeToString(raw)), 0600); err != nil {
		log.Fatalf("Failed to write nodekey to %s: %v", keyFile, err)
	}

	// Generate enode URL manually with IP and ports
	ip := net.ParseIP("127.0.0.1")
	tcpPort := 30303
	udpPort := 30303
	node = enode.NewV4(&privKey.PublicKey, ip, tcpPort, udpPort)

	fmt.Println("✅ Nodekey saved to:", keyFile)
	fmt.Println("🔗 Enode URL:", node.String())
	fmt.Println("🧪 Node ID:", node.ID().String())
	fmt.Println("📎 Public key:", hex.EncodeToString(crypto.FromECDSAPub(&privKey.PublicKey)))
	fmt.Println("🔐 Private key:", hex.EncodeToString(raw))
}
