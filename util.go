package main

func Union(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	for _, str := range arr2 {
		m[str] = true
	}
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func Intersect(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	res := make([]string, 0, len(m))
	for _, str := range arr2 {
		if m[str] {
			res = append(res, str)
		}
	}
	return res
}

func Difference(arr1 []string, arr2 []string) []string {
	m := make(map[string]bool)
	for _, str := range arr1 {
		m[str] = true
	}
	for _, str := range arr2 {
		delete(m, str)
	}
	res := make([]string, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
