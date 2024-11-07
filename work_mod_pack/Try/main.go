package main

import "work_mod_pack/package1/vcard"

func main() {
	inputs := []string{
		"Baiyun district Sanyuanli Aven No.1",
		"Huangshi Rd Baiyun district No.1313",
	}
	inputs_n := []string{
		"Kent",
		"Joden",
	}
	vcard.Address_input(inputs)
	vcard.Vcard_input(inputs_n)
	vcard.Vcard_output()
}
