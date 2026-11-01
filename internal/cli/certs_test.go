package cli

import "testing"

func TestCertGenerateECDSA(t *testing.T) {
	body := map[string]any{
		"name": "edge", "subject_cn": "edge.example", "key_algorithm": "ecdsa",
		"validity_days": float64(365), "ec_curve": "p-256", "hash_algorithm": "sha384",
		"subject_o": "CertLocker", "groups": []string{"ops"},
	}
	if err := normalizeCertGenerate(body); err != nil {
		t.Fatal(err)
	}
	if body["key_algorithm"] != "ECDSA" || body["ec_curve"] != "P-256" || body["hash_algorithm"] != "SHA384" {
		t.Fatalf("%v", body)
	}
	if _, ok := body["rsa_key_size"]; ok {
		t.Fatalf("rsa field present: %v", body)
	}
	if body["subject_cn"] != "edge.example" || body["subject_o"] != "CertLocker" {
		t.Fatalf("%v", body)
	}
	if body["validity_days"] != 365 {
		t.Fatalf("days: %v", body["validity_days"])
	}
}

func TestCertGenerateDefaults(t *testing.T) {
	body := map[string]any{"name": "edge", "subject_cn": "edge.example"}
	if err := normalizeCertGenerate(body); err != nil {
		t.Fatal(err)
	}
	if body["key_algorithm"] != "ECDSA" || body["ec_curve"] != "P-256" || body["hash_algorithm"] != "SHA256" || body["validity_days"] != 365 {
		t.Fatalf("%v", body)
	}
}

func TestCertGenerateMLDSAOmitsHash(t *testing.T) {
	body := map[string]any{
		"name": "pq", "subject_cn": "pq.example", "key_algorithm": "ML_DSA",
		"validity_days": 30, "ml_dsa_parameter_set": "ML_DSA_44", "hash_algorithm": "SHA512",
	}
	if err := normalizeCertGenerate(body); err != nil {
		t.Fatal(err)
	}
	if body["ml_dsa_parameter_set"] != "ML_DSA_44" {
		t.Fatalf("%v", body)
	}
	if _, ok := body["hash_algorithm"]; ok {
		t.Fatalf("hash sent for ML-DSA: %v", body)
	}
}

func TestCertGenerateRejectsBadDays(t *testing.T) {
	body := map[string]any{"name": "edge", "subject_cn": "edge.example", "key_algorithm": "RSA", "validity_days": 0}
	if err := normalizeCertGenerate(body); err == nil {
		t.Fatal("expected error")
	}
}

func TestCertGenerateRejectsBadHash(t *testing.T) {
	body := map[string]any{"name": "edge", "subject_cn": "edge.example", "hash_algorithm": "SHA1"}
	if err := normalizeCertGenerate(body); err == nil {
		t.Fatal("expected error")
	}
}
