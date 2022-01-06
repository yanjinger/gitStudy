package main

func main() {
	s := make([]int, 3, 10)
	_ = f(s)
}

func f(s []int) int {
	return s[1]
}
