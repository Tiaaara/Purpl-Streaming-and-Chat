package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: go run client/token_gen.go <policy_file.json>")
	}
	policyFile := os.Args[1]

	// baca isi file policy_all_denied.json
	data, err := os.ReadFile(policyFile)
	if err != nil {
		log.Fatalf("read policy file: %v", err)
	}

	var policy map[string]any
	if err := json.Unmarshal(data, &policy); err != nil {
		log.Fatalf("unmarshal: %v", err)
	}

	// ambil langsung ke dalam "policy"
	claims := jwt.MapClaims{
		"policy": policy["services"].(map[string]any)["trackingService-maximal"].(map[string]any)["purpose1"],
		"iss":    "tokenGenerator",
		"exp":    time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	key, err := os.ReadFile("client/private.pem")
	if err != nil {
		log.Fatalf("read key: %v", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		log.Fatalf("parse key: %v", err)
	}

	signed, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("sign token: %v", err)
	}

	fmt.Print(signed)
}
