package tax

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateTax(t *testing.T) {
	amount := 500.00
	expected := 5.0

	result, err := CalculateTax(amount)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCalculateTaxAndSave(t *testing.T) {
	repository := &TaxRepositoryMock{}
	repository.On("SaveTax", 10.0).Return(nil)
	repository.On("SaveTax", 0.0).Return(errors.New("error saving tax."))

	err := CalculateTaxAndSave(1000.0, repository)
	assert.Nil(t, err)

	err = CalculateTaxAndSave(0.0, repository)
	assert.NotNil(t, err, "error saving tax.")

	repository.AssertExpectations(t)
	repository.AssertCalled(t, "SaveTax", 10.0)

}
