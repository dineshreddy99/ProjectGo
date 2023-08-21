package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "yankay",
		Price: 1.56,
		SKU:   "abc-abc-cvcvcv",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
