package models

type Expr struct {
	Type     string  `json:"type"`               // "atom", "and", "or"
	Value    *string `json:"value,omitempty"`    // only for Type="atom"
	Operands []*Expr `json:"operands,omitempty"` // for Type="and"/"or"
}
