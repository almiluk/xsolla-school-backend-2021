package models

/// Response expressing occurred error and product with same SKU as the product in request in case of collision
type ResponseErrorProduct struct {
	//occurred error
	Error string `json:"error"`
	//product with same SKU as the product in request in case of collision
	Product Product `json:"product"`
} //@name ResponseErrorProduct

// Response expressing
type ResponseError struct {
	Error string `json:"error"` //occurred error
} // @name ResponseError
