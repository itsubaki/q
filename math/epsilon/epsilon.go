package epsilon

// E13 returns 1e-13 if eps is empty, otherwise returns eps[0].
func E13(eps ...float64) float64 {
	return Default(1e-13, eps...)
}

// E3 returns 1e-3 if eps is empty, otherwise returns eps[0].
func E3(eps ...float64) float64 {
	return Default(1e-3, eps...)
}

// Default returns v if eps is empty, otherwise returns eps[0].
func Default(v float64, eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return v
}
