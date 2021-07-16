package models

type Product struct {
	InputProduct
	Id int64
}

func NewProduct(SKU string, Name string, Type string, Cost uint, id int64) *Product {
	return &Product{InputProduct{SKU, Name, Type, Cost}, id}
}

func EmptyProduct() *Product {
	return &Product{*EmptyInputProduct(), 0}
}

type InputProduct struct {
	SKU  string
	Name string
	Type string
	Cost uint
}

func EmptyInputProduct() *InputProduct {
	return &InputProduct{"", "", "", 0}
}
