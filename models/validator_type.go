package models

type validate func() bool

type Validator interface {
	GetValidates() []validate
}

func Validate(v Validator) bool {
	for _, validate := range v.GetValidates() {
		if !validate() {
			return false
		}
	}
	return true
}
