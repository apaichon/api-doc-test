package converter

func MaskCreditCardNumer(creditCardNumber string) string {
	return creditCardNumber[:4] + "****" + creditCardNumber[len(creditCardNumber)-4:]
}

func MaskCreditCardExpiryDate(expiryDate string) string {
	return expiryDate[:2] + "/" + expiryDate[len(expiryDate)-2:]
}

func MaskCreditCardCVV(cvv string) string {
	return "**" + cvv[len(cvv)-1:]
}
