package data

import (
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator"
	_ "github.com/lib/pq"
	"io"
	"log"
	"regexp"
)

var DB *sql.DB

func init() {
	connStr := "user=postgres password=4360 dbname=postgres host=172.20.0.3 sslmode=disable" // Update with your actual connection string

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	/*if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	//defer DB.Close()
	log.Println("Connected to the database")*/
}

// Product defines the structure for an API product
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	//CreatedOn   string  `json:"-"`
	//UpdatedOn   string  `json:"-"`
	//DeletedOn   string  `json:"-"`
}

func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)

	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-absd-dfsdf
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

// Products is a collection of Product
type Products []*Product

// ToJSON serializes the contents of the collection to JSON
// NewEncoder provides better performance than json.Unmarshal as it does not
// have to buffer the output into an in-memory slice of bytes
// this reduces allocations and the overheads of the service
//
// https://golang.org/pkg/encoding/json/#NewEncoder
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func GetProducts() (Products, error) {
	rows, err := DB.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Products Products
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.SKU)
		if err != nil {
			return nil, err
		}
		Products = append(Products, &p)
	}
	return Products, nil
}

func AddProduct(p *Product) {
	query := "INSERT INTO products(name, description, price, sku) VALUES($1, $2, $3, $4) RETURNING id"
	err := DB.QueryRow(query, p.Name, p.Description, p.Price, p.SKU).Scan(&p.ID)
	if err != nil {
		log.Println("Error adding product:", err)
		return
	}
}

func UpdateProduct(id int, p *Product) error {
	query := "UPDATE products SET name=$1, description=$2, price=$3, sku=$4 WHERE id=$5"
	_, err := DB.Exec(query, p.Name, p.Description, p.Price, p.SKU, id)
	if err != nil {
		return err
	}
	return nil
}

//var ErrProductNotFound = fmt.Errorf("Product not found")
