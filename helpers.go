package shopping

func filter(vs []interface{}, f func(interface{}) bool) []interface{} {
	vsf := make([]interface{}, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
