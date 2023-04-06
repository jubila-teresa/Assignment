package models

import "time"

type Product struct {
	ProductName string          `json:"product_name"`
	Price       float32         `json:"price"`
	Category    ProductCategory `json:"category"`
	Quantity    int             `json:"quantity"`
}

type ProductResponse struct {
	ID          uint64          `json:"product_id"`
	ProductName string          `json:"product_name"`
	Price       float32         `json:"price"`
	Category    ProductCategory `json:"category"`
	Quantity    int             `json:"quantity"`
}

type Order struct {
	DispatchDate *time.Time  `json:"dispatch_date,omitempty"`
	Status       OrderStatus `json:"order_status"`
	Products     []Product   `json:"products"`
	Total        float32     `json:"total_amount"`
}

type OrderResponse struct {
	ID           uint64      `json:"order_id"`
	DispatchDate *time.Time  `json:"dispatch_date,omitempty"`
	Status       OrderStatus `json:"order_status"`
	Products     []Product   `json:"products"`
	Total        float32     `json:"total_amount"`
}

type OrderRequest struct {
	Products []struct {
		ProductID uint64 `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"products"`
}

type ProductCategory string

const (
	Premium ProductCategory = "Premium"
	Regular ProductCategory = "Regular"
	Budget  ProductCategory = "Budget"
)

func (p ProductCategory) IsValid() bool {
	switch p {
	case Premium, Regular, Budget:
		return true
	}
	return false
}

type OrderStatus string

const (
	Placed     OrderStatus = "Placed"
	Dispatched OrderStatus = "Dispatched"
	Completed  OrderStatus = "Completed"
	Cancelled  OrderStatus = "Cancelled"
)

func (o OrderStatus) IsValid() bool {
	switch o {
	case Placed, Dispatched, Completed, Cancelled:
		return true
	}
	return false
}

type UpdateOrder struct {
	DispatchDate string      `json:"dispatch_date"`
	Status       OrderStatus `json:"order_status"`
}
