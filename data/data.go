package data
import (
	"os"
	"log"
	"bufio"
	"fmt"
)

type Point struct {
	Label  string
	Values []float64
}

//func getLabel(label string, labels []string) string {
//	labels.g
//}

func getData(filename string, labels []string) []Point {
	// Open file
	file, err := os.Open("./spam/train.list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read from file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	data := make([]Point, 10)
	return data
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