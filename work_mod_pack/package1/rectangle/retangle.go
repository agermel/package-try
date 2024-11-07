package rectangle

type Rectangle struct {
	Width  int
	Length int
}

func Area(rec Rectangle) int {
	return rec.Width * rec.Length
}

func Perimeter(rec Rectangle) int {
	return (rec.Width + rec.Length) * 2
}
