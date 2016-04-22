package main
import (
	"testing"
	"github.com/grant/spamorham/data"
	"strconv"
	"fmt"
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

func TestLeftSplit(t *testing.T) {
	left, right := splitData(collegeData, 0, 25)
	for _, point := range left {
		if point.Values[0] >= 25 {
			t.Fail()
		}
		if len(left) != 3 {
			t.Fail()
		}
	}

	for _, point := range right {
		if point.Values[0] < 25 {
			t.Fail()
		}
		if len(right) != 7 {
			t.Fail()
		}
	}
}

func TestThreshold(t *testing.T) {
	gain, thresh := findBestThresHoldFast(collegeData, 1)
	if !floatEquals(gain, 0.321928094887) {
		t.Error(gain)
	}
	if thresh != 38000 {
		t.Error(thresh)
	}
}

func TestBestSplit(t *testing.T) {
	feature, thresh := findBestSplit(collegeData)
	if feature != 1 {
		t.Error(feature)
	}
	if thresh != 38000 {
		t.Error(thresh)
	}

	left, right := splitData(collegeData, feature, thresh)
	feature, thresh = findBestSplit(left)
	if feature != -1 {
		t.Error(feature)
	}
	if thresh != -1 {
		t.Error(thresh)
	}

	feature, thresh = findBestSplit(right)
	if feature != 0 {
		t.Error(feature)
	}
	if thresh != 43 {
		t.Error(thresh)
	}
}

func TestSubmission(t *testing.T) {
	train := data.GetSpamTrainData()
	valid := data.GetSpamValidData()
	preds := submission(train, valid)
	acc := accuracy(valid, preds)
	fmt.Println("Your current accuracy is: " + FloatToString(acc))
	if acc < .75 {
		t.Error(acc)
	}
}
