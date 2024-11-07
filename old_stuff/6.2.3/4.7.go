package main

func main() {
	a, b := 92, 13
	var ans int
	add(a, b, &ans)
	println(ans)
}

func add(a, b int, ans *int) {
	*ans = a * b //我说这是还原地址符
}

//笔记第十点的简单运用
