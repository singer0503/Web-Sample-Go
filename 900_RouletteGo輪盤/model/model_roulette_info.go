package model

type RoulletteInfo struct {
	// 數碼蛋的唯一識別碼，格式為uuid v4
	Id string `json:"id,omitempty"`
	// 數碼蛋的名稱
	Name string `json:"name,omitempty"`
	// 數碼蛋此時的狀態
	Status string `json:"status,omitempty"`
}
