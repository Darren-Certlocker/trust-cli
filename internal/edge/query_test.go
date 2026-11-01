package edge

import (
	"encoding/json"
	"testing"
)

func TestListQuery(t *testing.T) {
	raw, err := json.Marshal(ListQuery(ListFilter{
		Limit:  50,
		Skip:   10,
		Name:   "edge",
		Status: "ACTIVE",
		Group:  "ops",
	}))
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if body["limit"].(float64) != 50 || body["skip"].(float64) != 10 {
		t.Fatalf("%s", raw)
	}
	musts := body["musts"].([]any)
	if len(musts) != 2 {
		t.Fatalf("musts %s", raw)
	}
	ins := body["ins"].([]any)
	group := ins[0].(map[string]any)
	if group["field"] != "groups" {
		t.Fatalf("group %s", raw)
	}
}

func TestAuditQueryRange(t *testing.T) {
	raw, err := json.Marshal(AuditQuery(AuditFilter{
		Limit:   50,
		Action:  "HOST_ADD",
		Outcome: "SUCCESS",
		From:    "2026-11-01T00:00:00Z",
		To:      "2026-11-02T00:00:00Z",
		Email:   "ada@example.test",
	}))
	if err != nil {
		t.Fatal(err)
	}
	var body struct {
		Musts []Match `json:"musts"`
		Range []Range `json:"range_queries"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Musts) != 3 {
		t.Fatalf("musts %+v", body.Musts)
	}
	if len(body.Range) != 1 || body.Range[0].Field != "occurred_at" || body.Range[0].Type != "DATE" {
		t.Fatalf("range %+v", body.Range)
	}
	if body.Range[0].From != "2026-11-01T00:00:00Z" || body.Range[0].To != "2026-11-02T00:00:00Z" {
		t.Fatalf("range %+v", body.Range)
	}
}
