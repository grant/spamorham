package data
import (
	"os"
	"log"
	"bufio"
	"strings"
	"unicode"
)

type Point struct {
	Label  string
	Values []float64
}

func getLabel(ss []string, labels []string) string {
	for _, label := range labels {
		for _, s := range ss {
			if label == s {
				return label
			}
		}
	}
	return ""
}

func stringStatistics(s string) []float64 {
	s = strings.ToLower(s)
	values := make([]float64, 32)
	total := len(values)
	for _, c := range s {
		n := int(c)
		if 97 <= n && n <= 122 {
			values[n - 97] += 1
		} else if c == '.' {
			values[26] += 1
		} else if c == ',' {
			values[27] += 1
		} else if c == '?' {
			values[28] += 1
		} else if c == '!' {
			values[29] += 1
		} else if unicode.IsDigit(c) {
			values[30] += 1
		} else {
			values[31] += 1
		}
	}
	for i := 0; i < len(values); i += 1 {
		values[i] /= float64(total)
	}
	return values
}

func getFile(filename string) string {
	spamham, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer spamham.Close()
	spamhamscanner := bufio.NewScanner(spamham)
	f := ""
	for spamhamscanner.Scan() {
		f += spamhamscanner.Text()
	}
	return f
}

func getData(filename string, labels []string) []Point {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read from file
	scanner := bufio.NewScanner(file)
	data := make([]Point, 0)
	for scanner.Scan() {
		line := "../" + strings.TrimSpace(scanner.Text())
		label := getLabel(strings.Split(line, "."), labels)
		f := getFile(line)
		values := stringStatistics(f)
		data = append(data, Point{label, values})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func GetSpamTrainData() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/train.list", labels)
}

func GetSpamValidData() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/valid.list", labels)
}

func GetSpamTrainDataMany() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/train_many.list", labels)
}

func GetSpamValidDataMany() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/valid_many.list", labels)
}

func GetSpamTrainDataOne() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/train_one.list", labels)
}

func GetSpamValidDataOne() []Point {
	labels := []string{"ham", "spam"}
	return getData("../spam/valid_one.list", labels)
}

func GetCollegeData() []Point {
	return []Point{
		Point{"College", []float64{24, 40000}},
		Point{"No College", []float64{53, 52000}},
		Point{"No College", []float64{23, 25000}},
		Point{"College", []float64{25, 77000}},
		Point{"College", []float64{32, 48000}},
		Point{"College", []float64{52, 110000}},
		Point{"College", []float64{22, 38000}},
		Point{"No College", []float64{43, 44000}},
		Point{"No College", []float64{52, 27000}},
		Point{"College", []float64{48, 65000}},
	}
}