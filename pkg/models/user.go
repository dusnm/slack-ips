package models

import (
	"net/url"
	"strconv"
)

type (
	User struct {
		ID                string `json:"id"`
		Username          string `json:"username"`
		Name              string `json:"name"`
		BankAccountNumber string `json:"bank_account_number"`
		City              string `json:"city"`
		IPSString         string `json:"ips_string"`
	}
)

func (u User) QRCodeURL(host string, secure bool, amount float64) string {
	uri := new(url.URL)
	if secure {
		uri.Scheme = "https"
	} else {
		uri.Scheme = "http"
	}

	query := uri.Query()
	uri.Host = host
	uri.Path = "/image"

	query.Add("userId", u.ID)
	if amount > 0 {
		query.Add("amount", strconv.FormatFloat(amount, 'f', 2, 64))
	}

	uri.RawQuery = query.Encode()

	return uri.String()
}
