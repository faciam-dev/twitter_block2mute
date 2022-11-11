package common

func Contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// スライスの指定要素を削除する
func Remove[T any](arr []T, i int) []T {
	return arr[:i+copy(arr[i:], arr[i+1:])]
}
