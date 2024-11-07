package main

type shape interface {
	calc() float64
}

type rectangle struct {
	width  float64
	height float64
}

type circle struct {
	radius float64
}

func (r rectangle) calc() float64 {
	return r.height * r.width
}

func (r circle) calc() float64 {
	return r.radius * 3.14
}

func main() {

}
