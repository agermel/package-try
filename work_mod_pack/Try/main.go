package main

import "package1/vcard"

func main() {
	inputs := []string{
		"Baiyun district Sanyuanli Aven No.1",
		"Huangshi Rd Baiyun district No.1313",
	}
	vcard.Address_input(inputs)
	vcard.Address_output()
}
