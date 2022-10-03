package models

type WagerValidator struct {
	Wager *Wager
}

func (v *WagerValidator) ValidateTotalWagerValue() bool {
	return v.Wager.TotalWagerValue > 0
}
func (v *WagerValidator) ValidateOdds() bool {
	return v.Wager.Odds > 0
}
func (v *WagerValidator) ValidateSellingPercentage() bool {
	return v.Wager.SellingPercentage > 0 && v.Wager.SellingPercentage <= 100
}
func (v *WagerValidator) ValidateSellingPrice() bool {
	threshold := float64(v.Wager.TotalWagerValue) * (float64(v.Wager.SellingPercentage) / 100.0)
	return v.Wager.SellingPrice > threshold
}
func (v *WagerValidator) GetValidates() []validate {
	validates := []validate{
		v.ValidateTotalWagerValue,
		v.ValidateOdds,
		v.ValidateSellingPercentage,
		v.ValidateSellingPrice,
	}
	return validates
}
