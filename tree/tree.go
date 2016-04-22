package main
import (
	"github.com/grant/spamorham/data"
	"math"
)


type Tree struct {
	Leaf         bool
	Prediction   string
	FeatureIndex int
	Threshold    float64
	Left         *Tree
	Right        *Tree
}

type Prediction map[string]float64
type Data []data.Point

func (d Data) Len() int {
	return int(len(d))
}

func (t *Tree) predict(point data.Point) string {
	if t.Leaf {
		return t.Prediction
	}
	if point.Values[t.FeatureIndex] < t.Threshold {
		return t.Left.predict(point)
	} else {
		return t.Right.predict(point)
	}
}

func mostLikelyClass(prediction Prediction) string {
	biggestValue := math.SmallestNonzeroFloat64
	var biggestKey string
	for k, v := range prediction {
		if v > biggestValue {
			biggestKey = k
		}
	}
	return biggestKey
}

func accuracy(data Data, predictions []Prediction) float64 {
	var total float64 = 0
	var correct float64 = 0
	for i := range data {
		point := data[i]
		pred := predictions[i]
		total += 1
		guess := mostLikelyClass(pred)
		if guess == point.Label {
			correct += 1
		}
	}
	return correct / total
}

func splitData(data Data, featureId int, threshold float64) (Data, Data) {
	var left Data
	var right Data
	for _, datum := range data {
		if datum.Values[featureId] < threshold {
			left = append(left, datum)
		} else {
			right = append(right, datum)
		}
	}
	return left, right
}

func countLabels(data Data) map[string]int {
	count := make(map[string]int)
	for _, datum := range data {
		count[datum.Label] += 1
	}
	return count
}

func countsToEntropy(counts map[string]int, countTotal int) float64 {
	var entropy float64 = 0
	for _, v := range counts {
		var p float64 = float64(v) / float64(countTotal)
		if p > 0 {
			entropy += p * math.Log2(p)
		}
	}
	return -entropy
}

func getEntropy(data Data) float64 {
	counts := countLabels(data)
	if len(counts) == 0 {
		return 0
	} else {
		return countsToEntropy(counts, data.Len())
	}
}

func findBestThresHoldFast(data Data, featureId int) (float64, float64) {
	entropy := getEntropy(data)
	var bestGain float64 = 0
	var bestThreshold float64 = 0
	//	sort.Sort(ByFeatureValue(people))
	sortedData := data
	countTotal := len(sortedData)
	var left Data
	var right Data = sortedData
	rightCounts := countLabels(right)
	leftCounts := countLabels(left)
	leftTotal := len(left)
	rightTotal := len(right)

	var lastFeatureValue float64 = -1
	for _, point := range sortedData {
		if point.Values[featureId] != lastFeatureValue {
			var leftEntropy float64 = 0
			if leftTotal > 0 {
				leftEntropy = countsToEntropy(leftCounts, leftTotal)
			}
			var rightEntropy float64 = 0
			if rightTotal > 0 {
				rightEntropy = countsToEntropy(rightCounts, rightTotal)
			}
			var curr float64 = (leftEntropy * float64(leftTotal) + rightEntropy * float64(rightTotal)) / float64(countTotal)
			var gain float64 = entropy - curr
			if gain > bestGain {
				bestGain = gain
				bestThreshold = point.Values[featureId]
			}
		}

		// Split for next time
		leftCounts[point.Label] += 1
		leftTotal += 1
		rightCounts[point.Label] -= 1
		rightTotal -= 1
		lastFeatureValue = point.Values[featureId]
	}

	return bestGain, bestThreshold
}

func findBestSplit(data Data) (int, float64) {
	if len(data) < 2 {
		return -1, -1
	}
	bestFeatureId := -1
	var bestGain float64 = 0
	var bestThreshold float64 = 0
	for featureId := 0; featureId < len(data[0].Values); featureId++ {
		gain, threshold := findBestThresHoldFast(data, featureId)
		if gain > bestGain {
			bestGain = gain
			bestFeatureId = featureId
			bestThreshold = threshold
		}
	}
	return bestFeatureId, bestThreshold
}

func makeLeaf(data Data) *Tree {
	tree := &Tree{}
	counts := countLabels(data)
	prediction := make(map[string]float64)
	for k, v := range counts {
		prediction[k] = float64(v) / float64(len(data))
	}
	return tree
}

func c45(data Data, maxLevels int) *Tree {
	if maxLevels <= 0 {
		return makeLeaf(data)
	}
	if len(countLabels(data)) == 1 {
		return makeLeaf(data)
	}
	bestFeature, bestThreshold := findBestSplit(data)
	if bestFeature == -1 {
		return makeLeaf(data)
	}
	left, right := splitData(data, bestFeature, bestThreshold)
	tree := &Tree{
		Leaf: false,
		FeatureIndex:bestFeature,
		Threshold:bestThreshold,
		Left: c45(left, maxLevels - 1),
		Right: c45(right, maxLevels - 1),
	}
	return tree
}

func testSubmission(data Data, test Data, depth int) ([]string, *Tree) {
	tree := c45(data, depth)
	predictions := make([]string, 10)
	for _, point := range test {
		predictions = append(predictions, tree.predict(point))
	}
	return predictions, tree
}