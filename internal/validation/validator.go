package validation

import (
	"regexp"
	"strconv"
)

// CleanString removes all non-numeric characters from a string.
func CleanString(data string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(data, "")
}

func allDigitsEqual(data string) bool {
	for i := 1; i < len(data); i++ {
		if data[i] != data[0] {
			return false
		}
	}
	return true
}

func ValidateCPF(cpf string) string {
	cpf = CleanString(cpf)

	if len(cpf) != 11 || allDigitsEqual(cpf) {
		return "invalid"
	}

	sum := 0
	for i, weight := range []int{10, 9, 8, 7, 6, 5, 4, 3, 2} {
		num, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return "invalid"
		}
		sum += num * weight
	}
	digit1 := (sum * 10) % 11
	if digit1 == 10 {
		digit1 = 0
	}
	if strconv.Itoa(digit1) != string(cpf[9]) {
		return "invalid"
	}

	sum = 0
	for i, weight := range []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2} {
		num, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return "invalid"
		}
		sum += num * weight
	}
	digit2 := (sum * 10) % 11
	if digit2 == 10 {
		digit2 = 0
	}

	if strconv.Itoa(digit2) == string(cpf[10]) {
		return "valid"
	}
	return "invalid"
}

func ValidateCNPJ(cnpj string) string {
	cnpj = CleanString(cnpj)

	if len(cnpj) != 14 || allDigitsEqual(cnpj) {
		return "invalid"
	}

	digit1 := calculateCNPJDigit(cnpj[:12], 5)
	if int(cnpj[12]-'0') != digit1 {
		return "invalid"
	}

	digit2 := calculateCNPJDigit(cnpj[:13], 6)
	if int(cnpj[13]-'0') != digit2 {
		return "invalid"
	}

	return "valid"
}

// calculateCNPJDigit calculates the check digit for a CNPJ.
func calculateCNPJDigit(cnpj string, initialWeight int) int {
	weights := []int{}
	weight := initialWeight

	for i := 0; i < len(cnpj); i++ {
		weights = append(weights, weight)
		weight--
		if weight < 2 {
			weight = 9
		}
	}

	sum := 0
	for i, digit := range cnpj {
		num, err := strconv.Atoi(string(digit))
		if err != nil {
			return -1
		}
		sum += num * weights[i]
	}

	rest := sum % 11
	digit := 11 - rest
	if digit >= 10 {
		return 0
	}
	return digit
}
