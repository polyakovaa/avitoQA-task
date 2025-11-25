package task_2

type ItemRequest struct {
	SellerID   int        `json:"sellerID"`
	Name       string     `json:"name"`
	Price      int        `json:"price"`
	Statistics Statistics `json:"statistics"`
}

type ItemResponse struct {
	ID         string      `json:"id"`
	SellerID   int         `json:"sellerId"`
	Name       string      `json:"name"`
	Price      int         `json:"price"`
	Statistics *Statistics `json:"statistics"`
	CreatedAt  string      `json:"createdAt"`
}

type Statistics struct {
	Likes     int `json:"likes"`
	ViewCount int `json:"viewCount"`
	Contacts  int `json:"contacts"`
}
