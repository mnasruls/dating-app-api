package models

type TokenMetaData struct {
	Id     string `json:"id"`
	RtId   string `json:"rt_id"`
	Exp    int64  `json:"exp"`
	Verify bool   `json:"verify"`
}
