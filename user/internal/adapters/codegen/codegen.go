package codegen

import (
	"math/rand"
	"time"

	"github.com/lordvidex/errs"
)

type RandomCodeGenerator struct {
	cl          int
	letterBytes string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// New creates a new RandomCodeGenerator
// codeLen is the length of the code to generate
func New(codeLen int) (*RandomCodeGenerator, error) {
	if codeLen <= 5 {
		errs.B().Code(errs.InvalidArgument).Msg("Code length must be greater than 5")
	}
	return &RandomCodeGenerator{
		cl:          codeLen,
		letterBytes: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567889",
	}, nil
}

// GenerateCode generates a random code of length f.cl
// The code is a string of digits only
func (f *RandomCodeGenerator) GenerateCode() (string, error) {
	b := make([]byte, f.cl)
	for i := range b {
		b[i] = f.letterBytes[rand.Intn(len(f.letterBytes))]
	}
	return string(b), nil
}

// VerifyCode verifies if the code is valid
// A valid code is a string of digits of length f.cl
// Returns true if the code is valid, false otherwise
func (f *RandomCodeGenerator) VerifyCode(code string) bool {
	if code == "" || len(code) != f.cl {
		return false
	}
	for _, r := range code {
		isExpectedByte := (r >= '0' || r <= '9') || (r >= 'A' || r <= 'Z') || (r >= 'a' || r <= 'z')
		if !isExpectedByte {
			return false
		}
	}
	return true
}
