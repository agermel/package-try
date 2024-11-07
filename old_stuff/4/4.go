package main

import (
	"fmt"
)

func main() {
	var temp int
	a := "信管sjn"
	b := "计科sjn"
	c := "前端sjn"
	d := "后端sjn"
	e := "产品sjn"
	f := "设计sjn"
	aa := [3]int{75, 60, 50}
	bb := [3]int{80, 70, 55}
	cc := [3]int{85, 80, 60}
	dd := [3]int{90, 85, 65}
	ee := [3]int{95, 90, 70}
	ff := [3]int{80, 60, 92}
	fmt.Println("排序后的sjn们：")
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if aa[j] < aa[j+1] {
				temp = aa[j+1]
				aa[j+1] = aa[j]
				aa[j] = temp
			}
		}
	}
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if bb[j] < bb[j+1] {
				temp = bb[j+1]
				bb[j+1] = bb[j]
				bb[j] = temp
			}
		}
	}
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if cc[j] < cc[j+1] {
				temp = cc[j+1]
				cc[j+1] = cc[j]
				cc[j] = temp
			}
		}
	}
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if dd[j] < dd[j+1] {
				temp = dd[j+1]
				dd[j+1] = dd[j]
				dd[j] = temp
			}
		}
	}
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if ee[j] < ee[j+1] {
				temp = ee[j+1]
				ee[j+1] = ee[j]
				ee[j] = temp
			}
		}
	}
	for i := 3; i > 1; i-- {
		for j := 0; j < i-1; j++ {
			if ff[j] < ff[j+1] {
				temp = ff[j+1]
				ff[j+1] = ff[j]
				ff[j] = temp
			}
		}
	}
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", e, ee[0], ee[1], ee[2])
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", d, dd[0], dd[1], dd[2])
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", c, cc[0], cc[1], cc[2])
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", b, bb[0], bb[1], bb[2])
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", f, ff[0], ff[1], ff[2])
	fmt.Printf("%s:外貌相似度: %d, 声音相似度: %d, 性格相似度: %d\n", a, aa[0], aa[1], aa[2])
}
