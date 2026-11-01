package edge

type Match struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type In struct {
	Field  string   `json:"field"`
	Values []string `json:"values"`
}

type Range struct {
	Type  string `json:"type"`
	Field string `json:"field"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type Query struct {
	Limit        int     `json:"limit"`
	Skip         int     `json:"skip,omitempty"`
	Musts        []Match `json:"musts,omitempty"`
	Ins          []In    `json:"ins,omitempty"`
	RangeQueries []Range `json:"range_queries,omitempty"`
}

type ListFilter struct {
	Limit  int
	Skip   int
	Name   string
	Status string
	Group  string
}

func ListQuery(filter ListFilter) Query {
	query := Query{Limit: filter.Limit, Skip: filter.Skip}
	if filter.Name != "" {
		query.Musts = append(query.Musts, Match{Field: "name", Value: filter.Name})
	}
	if filter.Status != "" {
		query.Musts = append(query.Musts, Match{Field: "status", Value: filter.Status})
	}
	if filter.Group != "" {
		query.Ins = append(query.Ins, In{Field: "groups", Values: []string{filter.Group}})
	}
	return query
}

type AuditFilter struct {
	Limit     int
	Skip      int
	Action    string
	Outcome   string
	Severity  string
	Category  string
	Actor     string
	Email     string
	From      string
	To        string
	AssetID   string
	AssetType string
	Group     string
}

func AuditQuery(filter AuditFilter) Query {
	query := Query{Limit: filter.Limit, Skip: filter.Skip}
	addMust(&query, "action", filter.Action)
	addMust(&query, "outcome", filter.Outcome)
	addMust(&query, "severity", filter.Severity)
	addMust(&query, "category", filter.Category)
	addMust(&query, "actor_type", filter.Actor)
	addMust(&query, "user_email", filter.Email)
	addMust(&query, "asset_id", filter.AssetID)
	addMust(&query, "asset_type", filter.AssetType)
	if filter.Group != "" {
		query.Ins = append(query.Ins, In{Field: "groups", Values: []string{filter.Group}})
	}
	if filter.From != "" || filter.To != "" {
		query.RangeQueries = append(query.RangeQueries, Range{
			Type:  "DATE",
			Field: "occurred_at",
			From:  filter.From,
			To:    filter.To,
		})
	}
	return query
}

func addMust(query *Query, field, value string) {
	if value == "" {
		return
	}
	query.Musts = append(query.Musts, Match{Field: field, Value: value})
}
