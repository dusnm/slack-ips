package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/dusnm/slack-ips/pkg/utils"
)

var (
	ErrEmptyName                        = errors.New("name cannot be empty")
	ErrInvalidCharacterInName           = errors.New("name contains an invalid character")
	ErrEmptyBankAccountNumber           = errors.New("bank account number cannot be empty")
	ErrInvalidDashedBankAccountNumber   = errors.New("bank account number in dashed format is malformed")
	ErrBankAccountNumberNotNumeric      = errors.New("bank account number contains non numeric characters")
	ErrBankAccountNumberBankCodeInvalid = errors.New("bank account number contains an invalid bank code")
	ErrBankAccountNumberTooShort        = errors.New("bank account number is too short")
	ErrBankAccountNumberTooLong         = errors.New("bank account number is too long")
	ErrBankAccountNumberChecksumInvalid = errors.New("bank account number contains an invalid checksum")
	ErrEmptyCity                        = errors.New("city cannot be empty")
	ErrInvalidCharacterInCity           = errors.New("city contains an invalid character")
)

var (
	specialCharacterSet = map[rune]struct{}{
		'!': {}, '(': {}, '/': {}, '@': {}, '}': {}, ')': {}, ':': {}, '[': {}, '~': {},
		'#': {}, '*': {}, ';': {}, ']': {}, '„': {}, '$': {}, '+': {}, '<': {}, '^': {},
		'%': {}, ',': {}, '=': {}, '_': {}, '"': {}, '&': {}, '-': {}, '>': {}, '`': {},
		'\\': {}, '.': {}, '?': {}, '{': {}, '\'': {},
	}

	serbianCharacterSet = map[rune]struct{}{
		// Latin uppercase
		'A': {}, 'B': {}, 'C': {}, 'Č': {}, 'Ć': {}, 'D': {}, 'Đ': {}, 'E': {}, 'F': {}, 'G': {},
		'H': {}, 'I': {}, 'J': {}, 'K': {}, 'L': {}, 'M': {}, 'N': {}, 'O': {}, 'P': {}, 'R': {},
		'S': {}, 'Š': {}, 'T': {}, 'U': {}, 'V': {}, 'Z': {}, 'Ž': {},
		// Latin lowercase
		'a': {}, 'b': {}, 'c': {}, 'č': {}, 'ć': {}, 'd': {}, 'đ': {}, 'e': {}, 'f': {}, 'g': {},
		'h': {}, 'i': {}, 'j': {}, 'k': {}, 'l': {}, 'm': {}, 'n': {}, 'o': {}, 'p': {}, 'r': {},
		's': {}, 'š': {}, 't': {}, 'u': {}, 'v': {}, 'z': {}, 'ž': {},
		// Cyrillic uppercase
		'А': {}, 'Б': {}, 'В': {}, 'Г': {}, 'Д': {}, 'Ђ': {}, 'Е': {}, 'Ж': {}, 'З': {}, 'И': {},
		'Ј': {}, 'К': {}, 'Л': {}, 'Љ': {}, 'М': {}, 'Н': {}, 'Њ': {}, 'О': {}, 'П': {}, 'Р': {},
		'С': {}, 'Т': {}, 'Ћ': {}, 'У': {}, 'Ф': {}, 'Х': {}, 'Ц': {}, 'Ч': {}, 'Џ': {}, 'Ш': {},
		// Cyrillic lowercase
		'а': {}, 'б': {}, 'в': {}, 'г': {}, 'д': {}, 'ђ': {}, 'е': {}, 'ж': {}, 'з': {}, 'и': {},
		'ј': {}, 'к': {}, 'л': {}, 'љ': {}, 'м': {}, 'н': {}, 'њ': {}, 'о': {}, 'п': {}, 'р': {},
		'с': {}, 'т': {}, 'ћ': {}, 'у': {}, 'ф': {}, 'х': {}, 'ц': {}, 'ч': {}, 'џ': {}, 'ш': {},
		// Space
		' ': {},
	}

	fullCharacterSet = utils.MergeMaps(specialCharacterSet, serbianCharacterSet)
)

type (
	Init struct {
		Name              string
		BankAccountNumber string
		City              string
		UserID            string
		UserName          string
	}
)

// Format
//
// Must be called only after a successful
// validation! Returns the same structure
// with expanded Bank account number.
func (i Init) Format() Init {
	// Dashed format
	var (
		bankCode       string
		number         string
		checksumString string
	)

	// Look, I know I'm doing this work twice. But, to me,
	// clear separation of concerns is more important. Sue me.
	if strings.Contains(i.BankAccountNumber, "-") {
		// Dashed format
		parts := strings.Split(i.BankAccountNumber, "-")
		bankCode, number, checksumString = parts[0], parts[1], parts[2]
	} else {
		// Numeric format
		bankCode, number, checksumString = i.BankAccountNumber[:3], i.BankAccountNumber[3:len(i.BankAccountNumber)-2], i.BankAccountNumber[len(i.BankAccountNumber)-2:]
	}

	numberWithoutChecksum, checksum := expandAccountNumber(bankCode, number, checksumString)

	return Init{
		Name:              i.Name,
		BankAccountNumber: strconv.FormatInt(numberWithoutChecksum*100+checksum, 10),
		City:              i.City,
		UserID:            i.UserID,
		UserName:          i.UserName,
	}
}

func (i Init) Validate() error {
	validators := []func() error{
		i.validateName,
		i.validateBankAccountNumber,
		i.validateCity,
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

func (i Init) validateName() error {
	if len(i.Name) == 0 {
		return ErrEmptyName
	}

	for _, c := range i.Name {
		_, ok := fullCharacterSet[c]
		if !ok {
			return fmt.Errorf("%w: %c", ErrInvalidCharacterInName, c)
		}
	}

	return nil
}

func (i Init) validateBankAccountNumber() error {
	if len(i.BankAccountNumber) == 0 {
		return ErrEmptyBankAccountNumber
	}

	// Dashed format
	if strings.Contains(i.BankAccountNumber, "-") {
		return i.validateDashedBankAccountNumber()
	}

	// Numeric format
	return i.validateNumericBankAccountNumber()
}

func (i Init) validateDashedBankAccountNumber() error {
	parts := strings.Split(i.BankAccountNumber, "-")
	if len(parts) != 3 {
		return ErrInvalidDashedBankAccountNumber
	}

	for _, part := range parts {
		for _, c := range part {
			if !unicode.IsDigit(c) {
				return ErrBankAccountNumberNotNumeric
			}
		}
	}

	bankCode, number, checksumString := parts[0], parts[1], parts[2]

	if len(bankCode) != 3 {
		return ErrBankAccountNumberBankCodeInvalid
	}

	if len(number) < 1 {
		return ErrBankAccountNumberTooShort
	}

	if len(number) > 13 {
		return ErrBankAccountNumberTooLong
	}

	if len(checksumString) != 2 {
		return ErrBankAccountNumberChecksumInvalid
	}

	numberWithoutChecksum, checksum := expandAccountNumber(bankCode, number, checksumString)
	if !checkBankAccountNumberChecksum(numberWithoutChecksum, checksum) {
		return fmt.Errorf("%w: checksum check failed", ErrBankAccountNumberChecksumInvalid)
	}

	return nil
}

func (i Init) validateNumericBankAccountNumber() error {
	if len(i.BankAccountNumber) < 6 {
		return ErrBankAccountNumberTooShort
	}

	if len(i.BankAccountNumber) > 18 {
		return ErrBankAccountNumberTooLong
	}

	for _, c := range i.BankAccountNumber {
		if !unicode.IsDigit(c) {
			return ErrBankAccountNumberNotNumeric
		}
	}

	bankCode, number, checksumString := i.BankAccountNumber[:3], i.BankAccountNumber[3:len(i.BankAccountNumber)-2], i.BankAccountNumber[len(i.BankAccountNumber)-2:]
	numberWithoutChecksum, checksum := expandAccountNumber(bankCode, number, checksumString)
	if !checkBankAccountNumberChecksum(numberWithoutChecksum, checksum) {
		return fmt.Errorf("%w: checksum check failed", ErrBankAccountNumberChecksumInvalid)
	}

	return nil
}

func (i Init) validateCity() error {
	if len(i.City) == 0 {
		return ErrEmptyCity
	}

	for _, c := range i.City {
		_, ok := serbianCharacterSet[c]
		if !ok {
			return fmt.Errorf("%w: %c", ErrInvalidCharacterInCity, c)
		}
	}

	return nil
}

func expandAccountNumber(
	bankCode string,
	number string,
	checksumString string,
) (numberWithoutChecksum, checksum int64) {
	// Imagine needing a library for left padding 🥶
	if len(number) < 13 {
		b := strings.Builder{}
		b.Grow(13)
		for range 13 - len(number) {
			b.WriteString("0")
		}

		b.WriteString(number)
		number = b.String()
	}

	// I already checked everything to be numeric,
	// so I suppose this should work without any errors.
	// Watch me eat my words!
	numberWithoutChecksum, _ = strconv.ParseInt(bankCode+number, 10, 64)
	checksum, _ = strconv.ParseInt(checksumString, 10, 64)
	return
}

func checkBankAccountNumberChecksum(numberWithoutChecksum int64, checksum int64) bool {
	// The bank account number uses a simple modulo 97 operation,
	// as described in the ISO/IEC 7064 standard, to compute the checksum.
	return 98-numberWithoutChecksum*100%97 == checksum
}

// ToIPSString
//
// Formats the data into an IPS string.
// The format is described here: https://ips.nbs.rs/PDF/pdfPreporukeValidacija.pdf
func (i Init) ToIPSString() string {
	var (
		// The order of keys matters, so I'm using a helper slice here.
		keys   = []string{"K", "V", "C", "R", "N", "I", "SF", "S"}
		fields = map[string]string{
			"K":  "PR",
			"V":  "01",
			"C":  "1",
			"R":  i.BankAccountNumber,
			"N":  fmt.Sprintf("%s\r\n%s", i.Name, i.City),
			"I":  "RSD0,00",
			"SF": "289",
			"S":  "Transakcije po nalogu građana",
		}
	)

	builder := strings.Builder{}
	for _, key := range keys {
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(fields[key])
		builder.WriteString("|")
	}

	return strings.TrimSuffix(builder.String(), "|")
}
