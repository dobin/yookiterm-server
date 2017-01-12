package main

import (
		"fmt"
		"strings"
    "crypto/rsa"
    "crypto/x509"
		"crypto/sha1"
		"crypto"
    "encoding/pem"
    "io/ioutil"
    "log"
		"encoding/base64"
		"errors"
		"encoding/json"
)


func loadPrivateKey() (*rsa.PrivateKey, error) {
    // Read the private key
    pemData, err := ioutil.ReadFile("priv.key")
    if err != nil {
	return nil, err
        log.Fatalf("read key file: %s", err)
    }

    // Extract the PEM-encoded data block
    block, _ := pem.Decode(pemData)
    if block == nil {
        log.Fatalf("bad key data: %s", "not PEM-encoded")
		return nil, err
    }
    if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
        log.Fatalf("unknown key type %q, want %q", got, want)
		return nil, err
    }

    // Decode the RSA private key
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        log.Fatalf("bad private key: %s", err)
		return nil, err
    }

	return privateKey, nil
}


func loadHlPubKey() (*rsa.PublicKey, error) {
    pemData, err := ioutil.ReadFile("sso.hacking-lab.com.crt")
    if err != nil {
			return nil, err
        log.Fatalf("read key file: %s", err)
    }

    block, _ := pem.Decode(pemData)
    if block == nil {
        log.Fatalf("bad key data: %s", "not PEM-encoded")
	return nil, err
	}


    if got, want := block.Type, "CERTIFICATE"; got != want {
        log.Fatalf("unknown key type %q, want %q", got, want)
	return nil, err
    }

	var cert* x509.Certificate
	cert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)

	return rsaPublicKey, nil
}


func decryptToken(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Decrypt
	var result []byte
	n := 0
	blocksize := 512

	for n < len(data) {
    		var out []byte
        	out, err := rsa.DecryptPKCS1v15(nil, privateKey, data[n:n+blocksize])
	        if err != nil {
       	     		log.Fatalf("decrypt: %s", err)
			return nil, err
       		}
		n = n + blocksize
		result = append(result, out...)
	}

	return result, nil
}


func parseOuterToken(header string) ([]byte, error) {
	// Split token
	var tokenParts = strings.Split(header, ".")
	if len(tokenParts) != 2 {
		return nil, errors.New("Token wrong layout")
	}
	var content = tokenParts[1]

	// Decode
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return data, nil
}



func validateSignature(decryptedString string, publicKey *rsa.PublicKey) (bool, string, error) {
	parts := strings.Split(decryptedString, "::")
	if len(parts) != 2 {
		return false, "", errors.New("nah")
	}

	jsonStr := parts[0]
	signature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, jsonStr, err
	}

	h := sha1.New()
	h.Write([]byte(jsonStr))

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, h.Sum(nil), signature)
	if err != nil {
		return false, jsonStr, err
	}

	return true, jsonStr, nil
}



type HlJson struct {
	Citizen		string	`json:"citizen"`
	Admin		bool	`json:"admin"`
	Email		string	`json:"email"`
	Nickname	string	`json:"nickname"`
}


func getData(jsonStr string) (string, error) {
	fmt.Println("JSON: ", jsonStr)
	var hlJson HlJson

	if err := json.Unmarshal([]byte(jsonStr), &hlJson); err != nil {
		return "", err
	}

	return hlJson.Nickname, nil
}


func getUsername(token string) (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}

	hlPublicKey, err := loadHlPubKey()
	if err != nil {
		return "", err
	}

	data, err := parseOuterToken(token)
	if err != nil {
		return "", err
	}

	decrypted, err := decryptToken(data, privateKey)
	if err != nil {
		return "", err
	}

	decryptedString := string(decrypted[:])

	validSignature, jsonStr, err := validateSignature(decryptedString, hlPublicKey)
	if validSignature == false || err != nil {
		// yes
		//	return "", err
	}

	username, err := getData(jsonStr)
	if err != nil {
		return "", err
	}

	return username, nil
}
