package googletrans

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

func build_params(client string, query string, src string, dest string, token string) map[string]interface{} {
	params := map[string]interface{}{
		"client": client,
		"sl":     src,
		"tl":     dest,
		"hl":     dest,
		"dt":     []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		"ie":     "UTF-8",
		"oe":     "UTF-8",
		"otf":    1,
		"ssel":   0,
		"tsel":   0,
		"tk":     token,
		"q":      query,
	}
	return params
}
