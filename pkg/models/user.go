package models

import (
	"net"
	"net/url"
	"strconv"

	"github.com/dusnm/slack-ips/pkg/config"
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

func (u User) QRCodeURL(
	cfg config.App,
	amount float64,
) *url.URL {
	uri := new(url.URL)
	if cfg.Secure {
		uri.Scheme = "https"
	} else {
		uri.Scheme = "http"
	}

	query := uri.Query()
	if cfg.Port > 0 && !cfg.BehindProxy {
		uri.Host = net.JoinHostPort(cfg.Domain, strconv.FormatUint(uint64(cfg.Port), 10))
	} else {
		uri.Host = cfg.Domain
	}

	uri.Path = "/image"

	query.Add("userId", u.ID)
	if amount > 0 {
		query.Add("amount", strconv.FormatFloat(amount, 'f', 2, 64))
	}

	uri.RawQuery = query.Encode()

	return uri
}
