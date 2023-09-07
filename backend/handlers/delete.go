package handlers

import (
	"Practise/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		p.l.Println(err)
		return
	}
	p.l.Println("Deleting product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.DeleteProduct(id, &prod)
	if err != nil {
		p.l.Println(err)
		http.Error(rw, "error deleting product", http.StatusInternalServerError)
		return
	}

}
