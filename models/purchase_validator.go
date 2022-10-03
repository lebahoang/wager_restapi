package models

type PurchaseValidator struct {
	Purchase *Purchase
}

func (p *PurchaseValidator) ValidateBuyingPrice() bool {
	return p.Purchase.BuyingPrice > 0
}

func (p *PurchaseValidator) GetValidates() []validate {
	validates := []validate{p.ValidateBuyingPrice}
	return validates
}
