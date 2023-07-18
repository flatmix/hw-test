package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func writeAndClean(out *strings.Builder, write string, symbol *string) {
	if len(write) > 0 {
		out.WriteString(write)
	}
	*symbol = ""
}

func Unpack(inStr string) (string, error) {
	outStr := strings.Builder{}
	writeSym := ""
	backslash := false
	for _, sym := range inStr {
		num, errNum := strconv.Atoi(string(sym))
		if string(sym) == `\` {
			if !backslash {
				writeAndClean(&outStr, writeSym, &writeSym)
				backslash = true
			} else {
				writeSym = string(sym)
				backslash = false
			}
			continue
		}

		if errNum == nil && !backslash {
			if len(writeSym) == 0 {
				return "", ErrInvalidString
			}

			if num > 0 {
				writeAndClean(&outStr, strings.Repeat(writeSym, num), &writeSym)
			}
			writeSym = ""
		} else {
			if backslash && errNum != nil {
				return "", ErrInvalidString
			}
			writeAndClean(&outStr, writeSym, &writeSym)
			writeSym = string(sym)
			backslash = false
		}
	}
	writeAndClean(&outStr, writeSym, &writeSym)

	return outStr.String(), nil
}
