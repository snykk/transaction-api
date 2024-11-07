package requests

type WalletRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`  // price lebih besar dari 0
	Stock       int     `json:"stock" binding:"required,gte=0"` // stock tidak negatif

}

// func (productRequest *WalletRequest) ToDomain() *V1Domains.WalletDomain {
// 	return &V1Domains.WalletDomain{
// 		Name:        productRequest.Name,
// 		Description: productRequest.Description,
// 		Price:       productRequest.Price,
// 		Stock:       productRequest.Stock,
// 	}
// }
