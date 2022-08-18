package attestation

import (
	"os"
	"testing"
)

func TestAttestationParsing(t *testing.T) {
	content, _ := os.ReadFile("../../testdata/intoto_attestation.json")
	attPayload, _ := Parse(string(content))
	if attPayload.Subject == nil || attPayload.Predicate.Materials == nil {
		t.Error("Unable to parse attestation")
	}
}
