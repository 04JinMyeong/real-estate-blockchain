package vc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"realestate/crypto"
)

// GenerateAndSignVC issues a Verifiable Credential for a real estate agent.
// It validates the agent's qualification, constructs the VC JSON, and attaches a digital signature.
func GenerateAndSignVC(did, name, licenseNum, phone string) (string, error) {
    // 1. Validate broker qualification
    if !validateBroker(name, licenseNum) {
        return "", fmt.Errorf("broker validation failed for license: %s", licenseNum)
    }

    // 2. Read issuer DID from env
    issuerDID := os.Getenv("ISSUER_DID")
    if issuerDID == "" {
        return "", fmt.Errorf("ISSUER_DID environment variable is not set")
    }

    // 3. Build the credential document
    vc := map[string]interface{}{
        "@context":     []string{"https://www.w3.org/2018/credentials/v1"},
        "type":         []string{"VerifiableCredential", "RealEstateAgentVC"},
        "issuer":       issuerDID,
        "issuanceDate": time.Now().Format(time.RFC3339),
        "credentialSubject": map[string]string{
            "id":         did,
            "name":       name,
            "licenseNum": licenseNum,
            "phone":      phone,
            "status":     "valid",
        },
    }

    // 4. Sign the credentialSubject
    subjectBytes, err := json.Marshal(vc["credentialSubject"])
    if err != nil {
        return "", fmt.Errorf("failed to marshal credentialSubject: %w", err)
    }
    signature, err := crypto.Sign(subjectBytes)
    if err != nil {
        return "", fmt.Errorf("failed to sign VC: %w", err)
    }
    jws := base64.StdEncoding.EncodeToString(signature)

    // 5. Attach proof
    vc["proof"] = map[string]string{
        "type":               "Ed25519Signature2020",
        "created":            time.Now().Format(time.RFC3339),
        "proofPurpose":       "assertionMethod",
        "verificationMethod": issuerDID + "#key-1",
        "jws":                jws,
    }

    // 6. Return the final VC JSON
    vcBytes, err := json.MarshalIndent(vc, "", "  ")
    if err != nil {
        return "", fmt.Errorf("failed to marshal final VC: %w", err)
    }
    return string(vcBytes), nil
}
