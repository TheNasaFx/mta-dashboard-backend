package model

type PayMarketEbarimt struct {
	ID           int    `json:"id"`
	Pin          string `json:"pin"`
	CountReceipt int    `json:"count_receipt"` // Backward compatibility
	Cnt3         int    `json:"cnt_3"`         // Сүүлийн 3 хоногийн е-баримт
	Cnt30        int    `json:"cnt_30"`        // Сүүлийн 30 хоногийн е-баримт
	OpTypeName   string `json:"op_type_name"`  // Үйл ажиллагааны чиглэл
	MarName      string `json:"mar_name"`      // Зах нэр
	MarRegno     string `json:"mar_regno"`     // Зах регистр
	QrCode       string `json:"qr_code"`       // QR код
	MrchRegno    string `json:"mrch_regno"`    // Худалдаачийн регистр
}
