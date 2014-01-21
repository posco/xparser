package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

/*Parse control characters:
uppercase letter: absolute coordinates : 0
lowercase letter: relative coordinates : 1
not a control character: 2
*/
func IsControlCharacter(controlCharacter string) (controlValue int) {
	switch {
	case controlCharacter == "M":
		return 0
	case controlCharacter == "m":
		return 1
	case controlCharacter == "L":
		return 0
	case controlCharacter == "l":
		return 1
	}

	return 2
}

/* convert string coordinate pair to x, y values*/
func ConvertStringCoordinate(coordinateString string) (xValue float64, yValue float64) {

	xyStr := strings.Split(coordinateString, ",")
	xValue, xerr := strconv.ParseFloat(xyStr[0], 64) //parsed x values
	yValue, yerr := strconv.ParseFloat(xyStr[1], 64) //parsed x values

	if xerr != nil || yerr != nil {
		fmt.Println("Error pasing string")
		os.Exit(0)
	}

	return
}

/*Read coordinate metadata*/
func ReadCoordinate(fileName string) (splitString []string) {
	file, err := os.Open(fileName) // For read access.

	if err != nil {
		fmt.Println("Error opening file", file)
	}

	defer file.Close()

	data := make([]byte, 10000000)
	count, err := file.Read(data)
	if err != nil {
		fmt.Println("Error reading file", file)
	}

	readString := string(data[:count])

	splitString = strings.Split(readString, " ")

	return splitString
}

/*
M 31.9492,20.3438 l 26.4102,0 6.9297,0.1406
*/
func ParseStringCoordinate(readString []string) (x_axis []float64, y_axis []float64) {

	//fmt.Println(readString, len(readString))

	x_axis, y_axis = make([]float64, 0), make([]float64, 0)

	absFlag := true     //use absolute coordinate
	needOffset := false //offset is needed when m/l is encountered after M/L

	//setRef := false //set xRef, yRef when M/L is encounterred
	xRef, yRef := 0.0, 0.0
	//yRef := 0.0

	for _, str := range readString {

		if IsControlCharacter(str) == 0 {
			absFlag = true
			//fmt.Println("+++++++++", str)
		} else if IsControlCharacter(str) == 1 {
			absFlag = false
			needOffset = true
			//fmt.Println("*********", str)
		} else { //case 2:
			//fmt.Println("before", len(x_axis))
			x, y := ConvertStringCoordinate(str)
			//fmt.Println("---------", str)
			if absFlag { //absolute coordinate, no special treatment needed
				x_axis = append(x_axis, x)
				y_axis = append(y_axis, y)

				xRef, yRef = x, y
				//yRef = y
				absFlag = false

			} else {

				if needOffset {
					x += xRef
					y += yRef
					needOffset = false
				} else {
					//relative values
					x += x_axis[len(x_axis)-1]
					y += y_axis[len(y_axis)-1]
				}

				x_axis, y_axis = append(x_axis, x), append(y_axis, y)
				//y_axis = append(y_axis, y)
			}
			//fmt.Println("parsed value: ", x, " ,", y)
			//fmt.Println("after", len(x_axis))
		}

	}

	//fmt.Println(x_axis)
	//fmt.Println(y_axis)

	return
}

/*
Translate pixel coordinates to actual coordinates
Needed: reference points for lower left conor and upper right conor
*/
func TranslatePixelCoordinate(linear bool, startEndValues []float64, startEndPixels []float64, pixelsIn []float64) (actualCoordinates []float64) {

	actualCoordinates = make([]float64, 0)
	transCoordinate := 0.0

	for _, pixelValue := range pixelsIn {

		if linear {

			transCoordinate = (pixelValue-startEndPixels[0])*(startEndValues[1]-startEndValues[0])/
				(startEndPixels[1]-startEndPixels[0]) + startEndValues[0]

		} else { //logrithmic

			//their powers are linear
			coordinateRatio := (pixelValue-startEndPixels[0])*(math.Log10(startEndValues[1])-math.Log10(startEndValues[0]))/
				(startEndPixels[1]-startEndPixels[0]) + math.Log10(startEndValues[0])

			transCoordinate = math.Pow(10, coordinateRatio)

		}

		actualCoordinates = append(actualCoordinates, transCoordinate)
	}

	return
}

func WriteActualXY2File(fileName string, inValues []float64) {

	outFile, err := os.Create(fmt.Sprintf("%s.csv", fileName))

	if err != nil {
		fmt.Println("Error creating file", outFile)
	}

	defer outFile.Close()

	precision := 6

	for _, val := range inValues {
		fmt.Fprintf(outFile, "%s\n", strconv.FormatFloat(val, 'f', precision, 64))
	}

}

func main() {

	pixelXAxis, pixelYAxis := make([]float64, 0), make([]float64, 0)

	startXPixel, endXPixel, startXValue, endXValue := 31.9492, 1602.0892, 0.01, 100000000.0
	startYPixel, endYPixel, startYValue, endYValue := 20.0547, 435.6487, 0.0, 0.1

	keyfix := "pdf"

	inputString := ReadCoordinate(fmt.Sprintf("coordinate_%s.data", keyfix))
	pixelXAxis, pixelYAxis = ParseStringCoordinate(inputString)

	xAxis := TranslatePixelCoordinate(false, []float64{startXValue, endXValue}, []float64{startXPixel, endXPixel}, pixelXAxis)
	yAxis := TranslatePixelCoordinate(true, []float64{startYValue, endYValue}, []float64{startYPixel, endYPixel}, pixelYAxis)

	WriteActualXY2File("x", xAxis)
	WriteActualXY2File("y", yAxis)

}
