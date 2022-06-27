package epsilon

func E13(eps ...float64) float64 {
	return Default(1e-13, eps...)
}

func E3(eps ...float64) float64 {
	return Default(1e-3, eps...)
}

func Default(v float64, eps ...float64) float64 {
	if len(eps) > 0 {
		return eps[0]
	}

	return v
}
