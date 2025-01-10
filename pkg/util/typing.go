package util

func Float64RefFromInt(value int) *float64 {
	f := float64(value)
	return &f
}

func BoolTrueRef() *bool {
	b := true
	return &b
}

func BoolFalseRef() *bool {
	b := false
	return &b
}