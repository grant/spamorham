package main
import (
	"testing"
	"github.com/grant/spamorham/data"
	"strconv"
)

var collegeData Data = data.GetCollegeData()

// Helpers

var EPSILON float64 = 0.00000001
func floatEquals(a, b float64) bool {
	if ((a - b) < EPSILON && (b - a) < EPSILON) {
		return true
	}
	return false
}

func FloatToString(input float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input, 'f', 6, 64)
}


// Tests
func TestEntropy(t *testing.T) {
	entropy := getEntropy(collegeData)
	if !floatEquals(entropy, 0.9709505944546686) {
		t.Error("Wrong entropy " + FloatToString(entropy))
	}
}

