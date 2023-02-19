package codegen

import (
	"math/rand"
	"time"
)

type RandomCodeGenerator struct {
	cl int
}

// New creates a new RandomCodeGenerator
// codeLen is the length of the code to generate
func New(codeLen int) *RandomCodeGenerator {
	rand.Seed(time.Now().UnixNano())
	return &RandomCodeGenerator{cl: codeLen}
}

// GenerateCode generates a random code of length f.cl
// The code is a string of digits only
func (f *RandomCodeGenerator) GenerateCode() (string, error) {
	code := make([]byte, f.cl)
	for i := 0; i < f.cl; i++ {
		code[i] = byte(rand.Intn(10) + '0')
	}
	return string(code), nil
}

// VerifyCode verifies if the code is valid
// A valid code is a string of digits of length f.cl
// Returns true if the code is valid, false otherwise
func (f *RandomCodeGenerator) VerifyCode(code string) bool {
	if code == "" || len(code) != f.cl {
		return false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
