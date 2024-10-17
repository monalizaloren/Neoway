package validation

import (
	"regexp"
	"strconv"
)

func allDigitsEqual(data string) bool {
	for i := 1; i < len(data); i++ {
		if data[i] != data[0] {
			return false
		}
	}
	return true
}

func FormatCPF(cpf string) (string, bool) {
	re := regexp.MustCompile(`\D`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 || allDigitsEqual(cpf) {
		return "", false
	}

	formattedCPF := cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]

	if !ValidateCPF(cpf) {
		return "", false
	}
	return formattedCPF, true
}

func ValidateCPF(cpf string) bool {
	re := regexp.MustCompile(`\D`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 || allDigitsEqual(cpf) {
		return false
	}

	// First check digit calculation
	sum := 0
	for i, weight := range []int{10, 9, 8, 7, 6, 5, 4, 3, 2} {
		num, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return false
		}
		sum += num * weight
	}
	digit1 := (sum * 10) % 11
	if digit1 == 10 {
		digit1 = 0
	}
	if strconv.Itoa(digit1) != string(cpf[9]) {
		return false
	}

	// Second check digit calculation
	sum = 0
	for i, weight := range []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2} {
		num, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return false
		}
		sum += num * weight
	}
	digit2 := (sum * 10) % 11
	if digit2 == 10 {
		digit2 = 0
	}

	return strconv.Itoa(digit2) == string(cpf[10])
}

func FormatCNPJ(cnpj string) (string, bool) {
	re := regexp.MustCompile(`\D`)
	cnpj = re.ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 || allDigitsEqual(cnpj) {
		return "", false
	}

	formattedCNPJ := cnpj[:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:]

	if !ValidateCNPJ(cnpj) {
		return "", false
	}
	return formattedCNPJ, true
}

func ValidateCNPJ(cnpj string) bool {
	re := regexp.MustCompile(`\D`)
	cnpj = re.ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 || allDigitsEqual(cnpj) {
		return false
	}

	// First check digit validation
	digit1 := calculateCNPJDigit(cnpj[:12], 5)
	if int(cnpj[12]-'0') != digit1 {
		return false
	}

	// Second check digit validation
	//
	digit2 := calculateCNPJDigit(cnpj[:13], 6)
	if int(cnpj[13]-'0') != digit2 {
		return false
	}

	return true
}

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
