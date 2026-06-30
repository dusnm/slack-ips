package models

import "database/sql"

type Settings struct {
	QRFGColor  sql.NullString `json:"qr_fg_color"`
	QRBGColor  sql.NullString `json:"qr_bg_color"`
	QRShape    sql.NullString `json:"qr_shape"`
	QRCaption  sql.NullString `json:"qr_caption"`
	QRLogo     []byte         `json:"qr_logo"`
	QRShowLogo sql.NullBool   `json:"qr_show_logo"`
}

func (s Settings) GetQRFGColor() string {
	if s.QRFGColor.Valid {
		return s.QRFGColor.String
	}

	return "#000000"
}

func (s Settings) GetQRBGColor() string {
	if s.QRBGColor.Valid {
		return s.QRBGColor.String
	}

	return "#ffffff"
}

func (s Settings) GetQRShape() string {
	if s.QRShape.Valid {
		return s.QRShape.String
	}

	return "square"
}

func (s Settings) GetQRCaption() string {
	if s.QRCaption.Valid {
		return s.QRCaption.String
	}

	return ""
}

func (s Settings) ShouldShowLogo() bool {
	if s.QRShowLogo.Valid {
		return s.QRShowLogo.Bool
	}

	return false
}
