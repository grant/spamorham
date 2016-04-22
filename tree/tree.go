package main
import (
	"github.com/grant/spamorham/data"
	"math"
	"sort"
	"math/rand"
)


type Tree struct {
	Leaf         bool
	Prediction   Prediction
	FeatureIndex int
	Threshold    float64
	Left         *Tree
	Right        *Tree
}

type Prediction map[string]float64
type Data []data.Point

type SortedData struct {
	data      Data
	featureId int
}

func (p SortedData) Len() int { return len(p.data) }
func (p SortedData) Less(i, j int) bool { return p.data[i].Values[p.featureId] < p.data[j].Values[p.featureId]}
func (p SortedData) Swap(i, j int) { p.data[i], p.data[j] = p.data[j], p.data[i] }

func (d Data) Len() int {
	return int(len(d))
}

func (t *Tree) predict(point data.Point) Prediction {
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
			biggestValue = v
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
	sd := &SortedData{data, featureId}
	sort.Sort(sd)
	sortedData := sd.data
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
	var bestFeatureId int = -1
	var bestGain float64 = 0
	var bestThreshold float64 = -1
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
	tree := &Tree{
		Leaf:true,
	}
	counts := countLabels(data)
	prediction := make(Prediction)
	for k, v := range counts {
		prediction[k] = float64(v) / float64(len(data))
	}
	tree.Prediction = prediction
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

func testSubmission(data Data, test Data, depth int) ([]Prediction, *Tree) {
	tree := c45(data, depth)
	predictions := make([]Prediction, 0)
	for _, point := range test {
		predictions = append(predictions, tree.predict(point))
	}
	return predictions, tree
}

func submission(data Data, test Data) []Prediction {
	predictions, _ := testSubmission(data, test, 999)
	return predictions
}

// EC

func classifierFinal(data Data) []*Tree {
	depth := 10
	numTrees := 49
	trees := make([]*Tree, 0)
	for i := 0; i < numTrees; i += 1 {
		treeData := make(Data, 0)
		for j := 0; j < len(data); j += 1 {
			if rand.Intn(2) == 1 {
				treeData = append(treeData, data[j])
			}
		}
		tree := c45(treeData, depth)
		trees = append(trees, tree)
	}
	return trees
}

func predictFinal(model []*Tree, point data.Point) Prediction {
	predictionSum := make(map[string]float64)
	for i := 0; i < len(model); i += 1 {
		prediction := model[i].predict(point)
		for k, v := range prediction {
			predictionSum[k] += v
		}
	}

	// Normalize prediction sum
	var total float64 = 0
	for _, v := range predictionSum {
		total += v
	}
	for k, _ := range predictionSum {
		predictionSum[k] /= total
	}
	return predictionSum
}