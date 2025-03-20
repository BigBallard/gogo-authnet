package authnet

func Float64RefFromInt(value int) *float64 {
	f := float64(value)
	return &f
}

func Float64RefFromFloat(value float64) *float64 {
	return &value
}

func BoolTrueRef() *bool {
	b := true
	return &b
}

func BoolFalseRef() *bool {
	b := false
	return &b
}
