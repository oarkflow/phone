package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/oarkflow/pkg/str"
)

func main() {
	fmt.Println(os.Getwd())
	jsonFile, err := os.Open("./mcc-mnc-list.json")
	if err != nil {
		panic(err)
	}
	// read our opened jsonFile as a byte array.
	byteValue, _ := io.ReadAll(jsonFile)
	str.GenerateBinaryContent("phone", "networkMap", byteValue, "networks.go")

}

func verify() {
	start := time.Now()
	var phoneList []string
	str := `+84 8 825 6551 / 827 2601,848 820 2199,+84 8 910 7227,+84 8 910 6888,+84 8 847 6300,+84 8 829 8249,+1 408 1888 x 220,+84 8 914 0065,+84 8 843 9999,+65 6235 0077,+84 8 824 3562,+84 8 866 9308,+84 8 810 8789,+84 650 740 100,+84 8 829 9319,+84 8 822 4884,+84 613 514 190,+84 8 822 1199,08 8233597," +84 8 824 3833,"," +84 8 821 8632,"," +84 8 829 8430,"," +84 8 821 0121,"," +84 8 899 9375,",+84 8 821 9888,"  +84 8 744 4551,",+84 8 742 1604/ 6," +84 8 812 2196,"," +84 8 823 4999,"," +84 8 822 7737,"," +84 8 931 4700,"," +84 8 822 3203,"," +84 8 824 6125,"," +84 8 823 7993,"," +84 8 824 4115,"," +84 8 823 3412,",+84 8 824 2118,+84 8 822 6494,+84 8 825 7714," +84 8 829 9000,",+84 650 756 312,+84 8 822 7381," +84 8 824 3192,",+84 8 845 5528,+84 8 825 1723,+84 8 822 0566,+84 8 823 4326,+84 8 413 5050," +84 8 411 3333,",+84 8 930 2929," +84 8 821 9402,","  +84 8 837 4124,"," +84 8 437 1199,",+84 8 821 7841," +84 8 914 3747,"," +84 8 823 4151,"," +84 8 821 6141,"," +84 8 845 8518,",+84 8 829 2288,"  +84 8 824 1474,",+84 8 822 1612," +84 8 827 8888,"," +84 8 910 4855,"," +84 8 910 4855,"," +84 8 736 6200,"," +84 8 823 3372,"," +84 8 820 8496,",," +84 8 821 4812,"," +84 8 823 6466,"," +84 8 821 9266,",+84 8 822 2098," +84 8 824 3252,"," +84 8 827 3161,",+84 8 930 0394,+84 8 823 5266," +84 8 827 8008,"," +84 8 827 8000,",+84 8 821 9180," +84 8 829 8172,",+84 8 829 5368 / 822 6111," +84 8 829 5368,","  +84 8 822 6111 ext 1,"," +84 8 910 0492,","  +84 8 821 9437 / 8,"," +84 8 910 5120,"," +84 8 823 1333,"," +84 8 911 1508,"," +84 8 863 3634/ 862 0039,"," +84 8 411 9999,","  +1 800 577 755,"," +84 8 823 0796,"," +84 8 822 5678,"," +84 8 910 1660,"," +84 8 844 0611,"," +84 8 776 1349,"," +84 8 810 1888,"," +84 650 742 203,"," +84 8 824 6337,"," +84 8 822 0002,"," +84 8 413 0901,"," +84 8 823 9888,"," +84 8 865 7249,"," +84 8 829 8526,"," +84 8 836 8735,"," +84 8 823 9205,"," +84 8 822 9666,"," +84 8 824 6400,",+84 8 935 1920," +84 8 920 2020,",+84 8 822 8899," +84 8 829 8335/ 910 6898,"," +84 8 848 8888,",+84 8 825 7100,+84 650 742 137,+84 8 825 8980,"  +84 8 825 0490,"," +84 8 822 9433,",+84 8 413 5686," +84 8 824 2882,",+84 8 932 6598,+84 8 837 3031," +84 8 989 0390,"," +84 8 844 6667,"," +84 8 822 2755,"," +84 8 829 3649,"," +84 650 743 898,"," +84 8 833 6688,",," +84 58 820 090,",+84 8 827 5755," +84 8 925 0339,"," +84 8 823 9532,",`
	possiblePhones := strings.Split(str, `,`)
	re := regexp.MustCompile(`"([^"]*)"`)
	phones := re.FindAllString(str, -1)
	for _, element := range phones {
		element = strings.ReplaceAll(element, `"`, ``)
		element = strings.Trim(element, `,`)
		element = strings.TrimSpace(element)
		for i, ele := range possiblePhones {
			ele = strings.ReplaceAll(ele, `"`, ``)
			ele = strings.TrimSpace(ele)
			if strings.Contains(ele, element) {
				possiblePhones = RemoveIndex(possiblePhones, i)
				phoneList = append(phoneList, element)
			}
		}
	}
	for _, element := range possiblePhones {
		phoneList = append(phoneList, element)
	}
	phoneList = unique(phoneList)
	var p []string
	for i, element := range phoneList {
		element = strings.ReplaceAll(element, `"`, ``)
		element = strings.TrimSpace(element)
		if element == "" {
			phoneList = RemoveIndex(phoneList, i)
		} else {
			if strings.Contains(element, `/`) {
				s := strings.Split(element, `/`)
				s[0] = strings.TrimSpace(s[0])
				for t, str := range s {
					str = strings.TrimSpace(str)
					if t != 0 {
						rt := len(s[t]) - 1
						str = s[0][:len(s[0])-rt] + str
					}
					p = append(p, str)
				}
			} else {
				p = append(p, element)
			}
		}
	}
	fmt.Println(fmt.Sprintf("%s", time.Since(start)))
	/*
		fmt.Println(number)

		fmt.Println(phone.Verify("9856034616", "NP"))
		fmt.Println(phone.VerifyList([]string{"9856034616", "9856034617"}, "NP"))
		fmt.Println(phone.VerifyList([]string{"9856034616", "9856034617"}, "IN"))
		fmt.Println(phone.StatsByCarrier([]string{"9856034616", "9856034617", "919851446878"}, "NP"))
		fmt.Println(phone.StatsByCountry([]string{"9856034616", "9856034617", "919851446878"}, "NP"))
		fmt.Println(phone.Clean([]string{"9856034616", "9856034617", "919851446878"}, "NP"))*/
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
