package decimal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func numToWords(number int) string {
	var words string
	var unit = map[int]string{
		0: "ล้าน",
		1: "สิบ",
		2: "ร้อย",
		3: "พัน",
		4: "หมื่น",
		5: "แสน",
		6: "",
	}
	var value = map[int]string{
		0:  "",
		1:  "หนึ่ง",
		2:  "สอง",
		3:  "สาม",
		4:  "สี่",
		5:  "ห้า",
		6:  "หก",
		7:  "เจ็ด",
		8:  "แปด",
		9:  "เก้า",
		10: "เอ็ด",
		11: "ยี่",
	}
	lennumber := len(strconv.Itoa(number))
	for i := 0; i < lennumber; i++ {
		indexword := number % 10
		indexunit := i % 6
		if indexword == 0 {
			if indexunit == 0 {
				indexunit = 0
			} else {
				indexunit = 6
			}
		}
		if indexword == 1 {
			if indexunit == 0 {
				nextNumber := number / 10
				nextWord := nextNumber % 10
				if nextWord != 0 {
					indexword = 10
				} else {
					indexword = 1
				}
			}
			if i == lennumber-1 {
				indexword = 1
			}
			if indexunit == 1 {
				indexword = 0
			}
		}
		if indexword == 2 && indexunit == 1 {
			indexword = 11
		}
		if i == 0 {
			indexunit = 6
		}
		newword := value[indexword] + unit[indexunit]
		words = newword + words
		number = number / 10
	}
	return words
}

func fullNumToWords(s string) (string, error) {

	var decimal, number string

	if strings.Contains(s, ".") {
		split := strings.Split(s, ".")
		if len(split) != 2 {
			return "", errors.New(fmt.Sprintf("invalid (%s)  decimal format ", s))
		}
		number = split[0]
		decimal = split[1]
	} else {
		number = s
		decimal = "00"
	}

	i, err := strconv.Atoi(number)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot covert %s to int ", number))
	}

	i2, err := strconv.Atoi(decimal)
	if err != nil {
		return "", errors.New(fmt.Sprintf("cannot covert %s to int ", decimal))
	}

	if i == 0 && i2 == 0 {
		return "ศูนย์บาทถ้วน", nil
	}

	var decimalWord = ""
	if i2 == 0 {
		decimalWord = "ถ้วน"

	} else {
		decimalWord = "สตางค์"
	}

	if i == 0 && i2 > 0 {
		result := numToWords(i2) + decimalWord
		return result, nil
	}

	result := numToWords(i) + "บาท" + numToWords(i2) + decimalWord

	return result, nil
}
