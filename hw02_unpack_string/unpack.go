package main

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inStr string) (string, error) {

	outStr := strings.Builder{}
	formatStr := []rune(inStr)
	writeSym := ""
	ecran := false

	for _, sym := range formatStr {

		num, errNum := strconv.Atoi(string(sym))

		if string(sym) == `\` && !ecran {
			if len(writeSym) > 0 {
				outStr.WriteString(writeSym)
				writeSym = ""
			}
			ecran = true
			continue
		}

		if errNum == nil && !ecran {
			if len(writeSym) > 0 {
				if num > 0 {
					outStr.WriteString(strings.Repeat(writeSym, num))
				}
				writeSym = ""
			} else {
				return "", ErrInvalidString
			}
		} else {
			if ecran && string(sym) != `\` && errNum != nil {
				return "", ErrInvalidString
			}
			if len(writeSym) > 0 {
				outStr.WriteString(writeSym)
				writeSym = ""
			}

			if string(sym) != `\` || ecran {
				writeSym = string(sym)
				ecran = false
			}

		}
	}
	if len(writeSym) > 0 {
		outStr.WriteString(writeSym)
		writeSym = ""
	}

	return outStr.String(), nil
}
