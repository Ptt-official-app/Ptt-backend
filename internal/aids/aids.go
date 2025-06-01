package aids

import (
	"strconv"
	"strings"
)

type Aidu uint64

// Fn2Aidu 參考 pttbbs 裡面的 aids.c 將檔案名稱 M.123456789.A.0DC 的名稱轉換為 int32 的 Aidu 格式
func Fn2Aidu(fn string) Aidu {
	if fn == "" {
		return 0
	}

	var aidu Aidu
	var type_ Aidu
	var timestampPart Aidu
	var randomPart Aidu

	parts := strings.SplitN(fn, ".", 4)
	if len(parts) < 3 {
		return 0
	}

	switch parts[0] {
	case "M":
		type_ = 0
	case "G":
		type_ = 1
	default:
		return 0
	}

	v, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0
	}
	timestampPart = Aidu(v)
	if timestampPart == 0 {
		return 0
	}
	if parts[2] != "A" {
		return 0
	}

	if len(parts) < 4 || parts[3] == "" {
		randomPart = 0
	} else {
		v, err = strconv.ParseInt(parts[3], 16, 64)
		if err != nil {
			return 0
		}
		randomPart = Aidu(v)
	}
	if randomPart > 0xFFF {
		return 0
	}
	// Combine the parts into the Aidu format
	if timestampPart > 0xFFFFFFFF {
		return 0
	}
	if type_ > 0xF {
		return 0
	}

	aidu = (type_ << 44) | (timestampPart << 12) | randomPart

	return aidu
}

// Aidu2Aidc 參考 pttbbs 裡面的 aids.c 將原本的 Aidu 轉變為不同基數的編碼
func Aidu2Aidc(aidu Aidu) string {
	const aidu2aidcTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	const aidu2aidcTableSize = len(aidu2aidcTable)

	buf := make([]byte, 8)
	sp := len(buf) - 1

	for sp >= 0 {
		v := aidu % Aidu(aidu2aidcTableSize)
		aidu /= Aidu(aidu2aidcTableSize)
		buf[sp] = aidu2aidcTable[v]
		sp--
	}

	return string(buf[sp+1:])
}

// Aidu2Fn 參考 pttbbs 裡面的 aids.c 將 Aidu 轉換為檔案名稱 M.123456789.A.0DC 的格式
func Aidu2Fn(aidu Aidu) string {
	// Convert Aidu to M.123456789.A.0DC format
	type_ := (aidu >> 44) & 0xF
	timestamp := (aidu >> 12) & 0xFFFFFFFF
	random := aidu & 0xFFF

	var prefix string
	if type_ == 0 {
		prefix = "M"
	} else if type_ == 1 {
		prefix = "G"
	} else {
		return ""
	}

	return prefix + "." + strconv.FormatUint(uint64(timestamp), 10) + ".A." + strings.ToUpper(strconv.FormatUint(uint64(random), 16))
}

// Aidc2Aidu 參考 pttbbs 裡面的 aids.c 將不同基數的編碼轉換為 Aidu 格式
func Aidc2Aidu(aidc string) Aidu {
	var aidu Aidu

	for _, char := range aidc {
		var v Aidu
		switch {
		case char >= '0' && char <= '9':
			v = Aidu(char - '0')
		case char >= 'A' && char <= 'Z':
			v = Aidu(char - 'A' + 10)
		case char >= 'a' && char <= 'z':
			v = Aidu(char - 'a' + 36)
		case char == '-':
			v = 62
		case char == '_':
			v = 63
		default:
			return 0
		}
		aidu = (aidu << 6) | (v & 0x3F)
	}

	return aidu
}
