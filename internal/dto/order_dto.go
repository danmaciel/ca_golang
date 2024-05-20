package dto

type OrderInputDTO struct {
	Price float32 `json:"price"`
	Tax   float32 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float32 `json:"price"`
	Tax        float32 `json:"tax"`
	FinalPrice float32 `json:"final_price"`
}

type ListOrderOutputDTO struct {
	Orders []OrderOutputDTO `json:"orders"`
}
