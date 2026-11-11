package entity

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenAnEmptyID_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestGivenAnEmptyPrice_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "123"}
	assert.Error(t, order.IsValid(), "invalid price")
}

func TestGivenAnEmptyTax_WhenCreateANewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "123", Price: 10}
	assert.Error(t, order.IsValid(), "invalid tax")
}

func TestGivenAValidParams_WhenICallNewOrder_ThenIShouldReceiveCreateOrderWithAllParams(t *testing.T) {
	order := Order{
		ID:    "123",
		Price: 10.0,
		Tax:   2.0,
	}
	assert.Equal(t, "123", order.ID)
	assert.Equal(t, float32(10), order.Price)
	assert.Equal(t, float32(2), order.Tax)
	assert.Nil(t, order.IsValid())
}

func TestGivenAValidParams_WhenICallNewOrderFunc_ThenIShouldReceiveCreateOrderWithAllParams(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	assert.Equal(t, "123", order.ID)
	assert.Equal(t, float32(10), order.Price)
	assert.Equal(t, float32(2), order.Tax)
}

func TestGivenAPriceAndTax_WhenICallCalculatePrice_ThenIShouldSetFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	assert.Nil(t, order.CalculateFinalPrice())
	assert.Equal(t, float32(12), order.FinalPrice)
}

func TestOrderRejectsNonFinitePriceAndTax(t *testing.T) {
	cases := []struct {
		name  string
		price float32
		tax   float32
	}{
		{"NaN price", float32(math.NaN()), 1},
		{"positive infinite price", float32(math.Inf(1)), 1},
		{"negative infinite price", float32(math.Inf(-1)), 1},
		{"NaN tax", 10, float32(math.NaN())},
		{"positive infinite tax", 10, float32(math.Inf(1))},
		{"negative infinite tax", 10, float32(math.Inf(-1))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			order := &Order{ID: "123", Price: tc.price, Tax: tc.tax, FinalPrice: 77}
			assert.Error(t, order.IsValid())
			assert.Error(t, order.CalculateFinalPrice())
			assert.Equal(t, float32(77), order.FinalPrice,
				"invalid calculation must not mutate existing final price")
		})
	}
}

func TestOrderRejectsOverflowWithoutChangingFinalPrice(t *testing.T) {
	order := &Order{
		ID:         "123",
		Price:      math.MaxFloat32,
		Tax:        math.MaxFloat32,
		FinalPrice: 77,
	}
	assert.NoError(t, order.IsValid(), "inputs are finite and positive")
	assert.Error(t, order.CalculateFinalPrice(), "sum overflows float32")
	assert.Equal(t, float32(77), order.FinalPrice)
}

func TestOrderInvalidInputDoesNotChangeFinalPrice(t *testing.T) {
	order := &Order{ID: "123", Price: -1, Tax: 2, FinalPrice: 33}
	assert.Error(t, order.CalculateFinalPrice())
	assert.Equal(t, float32(33), order.FinalPrice)
}
