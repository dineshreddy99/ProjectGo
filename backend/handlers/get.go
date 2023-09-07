package handlers

import (
	"Practise/data"
	"net/http"
)

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	lp, err := data.GetProducts()
	if err != nil {
		p.l.Println(err)
		http.Error(rw, "Unable to get products", http.StatusInternalServerError)
		return
	}

	err = lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
