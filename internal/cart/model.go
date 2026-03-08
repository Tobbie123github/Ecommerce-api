package cart

type CartItem struct {
	ProductID string   `json:"productid"`
	Name      string   `json:"name"`
	Price     float64  `json:"price"`
	ImageURL  []string `json:"imageUrl"`
	Quantity  int      `json:"quantity"`
	Size      string   `json:"size"`
}

type Cart struct {
	UserId string     `json:"userId"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}
