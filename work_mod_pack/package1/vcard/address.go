package vcard

import (
	"fmt"
	"regexp"
	"strconv"
)

type Vcard struct {
	name    string
	Address *Address
}

type Address struct {
	street  string
	number  int
	distric string
}

var addresses []Address
var vcards []Vcard

//example1: Baiyun district Sanyuanli Aven No.1
//example2: Huangshi Rd Baiyun district No.1313

func Address_input(strs []string) {
	patdis := "([a-zA-Z]+)\\s+district\\s+"
	patstr := "([a-zA-Z]+)\\s+(Aven|Rd)+"
	patno := "No\\.+(\\d+)"

	redis, _ := regexp.Compile(patdis)
	restr, _ := regexp.Compile(patstr)
	reno, _ := regexp.Compile(patno)

	for _, str := range strs {
		dismatch := redis.FindStringSubmatch(str)
		strmatch := restr.FindStringSubmatch(str)
		nomatch := reno.FindStringSubmatch(str)

		number, _ := strconv.Atoi(nomatch[1])

		address := Address{
			distric: dismatch[1],
			street:  strmatch[1],
			number:  number,
		}
		addresses = append(addresses, address)
	}
}

func Address_output() {
	for _, address := range addresses {
		fmt.Printf("District:%s Street:%s No:%d \n",
			address.distric, address.street, address.number)
	}
}

func Vcard_input(names []string) {
	for i, address := range addresses {
		vcard := Vcard{
			name:    names[i],
			Address: &address,
		}
		vcards = append(vcards, vcard)
	}
}

func Vcard_output() {
	for _, vcard := range vcards {
		fmt.Printf("Name:%s ", vcard.name)
		fmt.Printf("District:%s Street:%s No:%d \n",
			vcard.Address.distric, vcard.Address.street, vcard.Address.number)
	}
}
