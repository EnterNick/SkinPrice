package skins

import (
	"fmt"
	"strings"
)

const LisSkinsCurrency = "1"

func NormalizeCurrencyCode(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "1", "USD":
		return "1"
	case "3", "EUR":
		return "3"
	case "5", "RUB":
		return "5"
	default:
		return "1"
	}
}

func FormatPriceText(priceCents int64, currency string) string {
	sign := ""
	if priceCents < 0 {
		sign = "-"
		priceCents = -priceCents
	}

	whole := priceCents / 100
	fraction := priceCents % 100

	switch currency {
	case "3":
		return fmt.Sprintf("%s€%d.%02d", sign, whole, fraction)
	case "5":
		if fraction == 0 {
			return fmt.Sprintf("%s%d ₽", sign, whole)
		}
		return fmt.Sprintf("%s%d.%02d ₽", sign, whole, fraction)
	default:
		return fmt.Sprintf("%s$%d.%02d", sign, whole, fraction)
	}
}

func NormalizePriceText(price *NewSkin, currency string) string {
	if price == nil {
		return ""
	}
	if price.PriceCents != nil {
		return FormatPriceText(*price.PriceCents, currency)
	}
	return price.PriceText
}

func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
