package entity

import (
	"errors"
	"math"
)

type Order struct {
	ID         string
	Price      float32
	Tax        float32
	FinalPrice float32
}

func NewOrder(id string, price float32, tax float32) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	err := order.IsValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) IsValid() error {
	if o.ID == "" {
		return errors.New("invalid id")
	}
	if math.IsNaN(float64(o.Price)) || math.IsInf(float64(o.Price), 0) || o.Price <= 0 {
		return errors.New("invalid price")
	}
	if math.IsNaN(float64(o.Tax)) || math.IsInf(float64(o.Tax), 0) || o.Tax <= 0 {
		return errors.New("invalid tax")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() error {
	if err := o.IsValid(); err != nil {
		return err
	}
	finalPrice := o.Price + o.Tax
	if math.IsNaN(float64(finalPrice)) || math.IsInf(float64(finalPrice), 0) {
		return errors.New("invalid final price")
	}
	o.FinalPrice = finalPrice
	return nil
}
