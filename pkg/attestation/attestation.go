package attestation

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

func Parse(attJson string) (*Payload, error) {
	att := Attestation{}
	err := json.Unmarshal([]byte(attJson), &att)
	if err != nil {
		return nil, err
	}
	decodedPayloadString, err := base64.StdEncoding.DecodeString(att.PayloadB64)
	if err != nil {
		return nil, err
	}
	payload := Payload{}
	err = json.Unmarshal(decodedPayloadString, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

type Attestation struct {
	PayloadB64 string `json:"payload"`
}

type Payload struct {
	Subject   []Subject `json:"subject"`
	Predicate Predicate `json:"predicate"`
}

type Subject struct {
	Name   string `json:"name"`
	Digest Digest `json:"digest"`
}

type Digest struct {
	Sha256 string `json:"sha256"`
}

type Material struct {
	Uri    string `json:"uri"`
	Digest Digest `json:"digest"`
}

type Predicate struct {
	Metadata  Metadata   `json:"metadata"`
	Materials []Material `json:"materials"`
}

type Metadata struct {
	BuildInvocationID string    `json:"buildInvocationID"`
	BuildStartedOn    time.Time `json:"buildStartedOn"`
}
