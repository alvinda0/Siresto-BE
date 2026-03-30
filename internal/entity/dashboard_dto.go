package entity

// DashboardHomeResponse adalah response untuk endpoint home/dashboard
type DashboardHomeResponse struct {
	TotalItemsByDate   []DailyStats             `json:"total_items_by_date"`  // Total item per tanggal (7 hari)
	RevenuByDate       []DailyStats             `json:"revenue_by_date"`      // Pendapatan per tanggal (7 hari)
	BestSellingDaily   []BestSellingItem        `json:"best_selling_daily"`   // Item terlaris hari ini
	BestSellingWeekly  []BestSellingItem        `json:"best_selling_weekly"`  // Item terlaris minggu ini
	BestSellingMonthly []BestSellingItem        `json:"best_selling_monthly"` // Item terlaris bulan ini
	ComplimentaryItems []ComplimentaryItemStats `json:"complimentary_items"`  // Item complimentary
}

// DailyStats adalah statistik per hari
type DailyStats struct {
	Date  string  `json:"date"`  // Format: 2024-01-15
	Value float64 `json:"value"` // Bisa total items atau revenue
}

// BestSellingItem adalah item terlaris
type BestSellingItem struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalQty    int64   `json:"total_qty"`
	TotalAmount float64 `json:"total_amount"`
}

// ComplimentaryItemStats adalah statistik item complimentary
type ComplimentaryItemStats struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalQty    int64  `json:"total_qty"`
}
