package bech32

import (
	"strings"
	"testing"
)

func TestBech32(t *testing.T) {
	tests := []struct {
		str   string
		valid bool
	}{
		{"A12UEL5L", true},
		{"a12uel5l", true},
		{"an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1tt5tgs", true},
		{"abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw", true},
		{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", true},
		{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", true},

		// invalid checksum
		{"split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w", false},
		// invalid character (space) in hrp
		{"s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p", false},
		{"split1cheo2y9e2w", false}, // invalid character (o) in data part
		{"split1a2y9w", false},      // too short data part
		{"1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false}, // empty hrp
		// invalid character (DEL) in hrp
		{"spl" + string(rune(127)) + "t1checkupstagehandshakeupstreamerranterredcaperred2y9e3w", false},
		// too long
		{"11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j", false},

		// BIP 173 invalid vectors.
		{"an84characterslonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio1569pvx", false},
		{"pzry9x0s0muk", false},
		{"1pzry9x0s0muk", false},
		{"x1b4n0q5v", false},
		{"li1dgmt3", false},
		{"de1lg7wt\xff", false},
		{"A1G7SGD8", false},
		{"10a06t8", false},
		{"1qzzfhee", false},
	}

	for _, test := range tests {
		str := test.str
		hrp, decoded, err := Decode(str)
		if !test.valid {
			// Invalid string decoding should result in error.
			if err == nil {
				t.Errorf("expected decoding to fail for invalid string %v", test.str)
			}
			continue
		}

		// Valid string decoding should result in no error.
		if err != nil {
			t.Errorf("expected string to be valid bech32: %v", err)
		}

		// Check that it encodes to the same string.
		encoded, err := Encode(hrp, decoded)
		if err != nil {
			t.Errorf("encoding failed: %v", err)
		}
		if encoded != str {
			t.Errorf("expected data to encode to %v, but got %v", str, encoded)
		}

		// Flip a bit in the string an make sure it is caught.
		pos := strings.LastIndexAny(str, "1")
		flipped := str[:pos+1] + string((str[pos+1] ^ 1)) + str[pos+2:]
		if _, _, err = Decode(flipped); err == nil {
			t.Error("expected decoding to fail")
		}
	}
}
