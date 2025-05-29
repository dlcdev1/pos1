package entity

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProduct(t *testing.T) {
	p, error := NewProduct("test", 100.0)
	assert.Nil(t, error)
	assert.NotNil(t, p)
	assert.NotEmpty(t, p.ID)
	assert.NotEmpty(t, p.Name)
	assert.NotEmpty(t, p.Price)
	assert.Equal(t, 100.0, p.Price)
	assert.Equal(t, "test", p.Name)
}

func TestProductWhenNameisRequired(t *testing.T) {
	p, error := NewProduct("", 100.0)
	assert.NotNil(t, error)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequired, error)
}

func TestProductWhenPriceIsRequired(t *testing.T) {
	p, error := NewProduct("test", 0)
	assert.NotNil(t, error)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, error)
}

func TestProductWhenPriceIsInvalid(t *testing.T) {
	p, error := NewProduct("test", -100)
	assert.NotNil(t, error)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, error)
}

func TestProductValidate(t *testing.T) {
	p, error := NewProduct("test", 100)
	assert.Nil(t, error)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
