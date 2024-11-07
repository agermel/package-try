package main

import (
	"fmt"
	"work_mod_pack/package1/rectangle"
)

func main() {
	//vcard
	/*inputs := []string{
		"Baiyun district Sanyuanli Aven No.1",
		"Huangshi Rd Baiyun district No.1313",
	}
	inputs_n := []string{
		"Kent",
		"Joden",
	}
	vcard.Address_input(inputs)
	vcard.Vcard_input(inputs_n)
	vcard.Vcard_output()*/

	//struct in pkg

	rec1 := rectangle.Rectangle{Width: 1, Length: 2}
	fmt.Println(rectangle.Area(rec1), rectangle.Perimeter(rec1))
}
