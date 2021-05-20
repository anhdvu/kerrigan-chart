package data

type WsMsg struct {
	M string `json:"m"`
	D struct {
		S string  `json:"s"`
		T int64   `json:"t"`
		V float64 `json:"v"`
		E float64 `json:"e"`
	} `json:"d"`
}
