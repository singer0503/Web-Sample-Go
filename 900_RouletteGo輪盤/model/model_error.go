package model

type ModelError struct {
	// 錯誤訊息
	Message string `json:"message"`
	// 錯誤代碼:  * `3000` - Internal error
	Code float64 `json:"code"`
}
