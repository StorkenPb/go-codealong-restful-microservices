package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "name",
		Price: 9999,
		SKU: "tes-ting-asd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}