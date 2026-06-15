package models

import "encoding/base64"

type QR struct {
	buff []byte
}

func NewQR(buff []byte) QR {
	return QR{buff: buff}
}

func (q QR) Base64Encode() string {
	return base64.StdEncoding.EncodeToString(q.buff)
}

func (q QR) Bytes() []byte {
	return q.buff
}
