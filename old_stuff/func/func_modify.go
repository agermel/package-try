package main

func main() {
	x := 10
	modify(&x)
	println(x)
}

func modify(x *int) {
	*x = 100
}
