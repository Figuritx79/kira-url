package base62

import (
	"math/big"
	"testing"
)

func TestEncodeToBase62(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  string
	}{
		{name: "nil input encodes to empty string", value: nil, want: ""},
		{name: "empty input encodes to empty string", value: []byte{}, want: ""},
		{name: "zero byte encodes to empty string", value: []byte{0}, want: ""},
		{name: "one encodes to first symbol", value: []byte{0x01}, want: "1"},
		{name: "last symbol encodes to z", value: []byte{61}, want: "z"},
		{name: "base rollover encodes to 10", value: []byte{62}, want: "10"},
		{name: "multi byte number encodes big endian", value: []byte{0x01, 0x00}, want: "48"},
		{name: "one million encodes to 4C92", value: []byte{0x0F, 0x42, 0x40}, want: "4C92"},
		{name: "max six symbol value encodes to zzzzzz", value: new(big.Int).SetInt64(56800235583).Bytes(), want: "zzzzzz"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := EncodeToBase62(test.value)

			if got != test.want {
				t.Fatalf("EncodeToBase62(%v) = %q, want %q", test.value, got, test.want)
			}
		})
	}
}

func TestRandamBase62Number(t *testing.T) {
	limit := big.NewInt(possibilities)

	for range 100 {
		number := RandamBase62Number()

		if number == nil {
			t.Fatal("RandamBase62Number() = nil, want a random number")
		}
		if len(number) > 5 {
			t.Fatalf("RandamBase62Number() length = %d, want at most 5 bytes for a value below 62^6", len(number))
		}
		if value := new(big.Int).SetBytes(number); value.Cmp(limit) >= 0 {
			t.Fatalf("RandamBase62Number() = %v, want a value below %v", value, limit)
		}
	}
}
