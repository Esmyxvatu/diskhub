package main

import (
	"math"
	"fmt"
)

//============================================================ Functions ============================================================================

func generatePieChart(data map[string]float64, width, height int) string {
	// Calculate the sum of the data
	var total float64
	for _, v := range data {
		total += v
	}

	// Create the main SVG tag
	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, 
		width + 5, height, width, height)

	// Calc part by part the angle of every part
	var startAngle float64
	for label, value := range data {
		percentage := value / total
		endAngle := startAngle + percentage*360

		if percentage > 0 {
			// If the percentage is between 0.5 and 1, round the percentage to 1
			largeArcFlag := 0
			if percentage > 0.5 {
				largeArcFlag = 1
			}

			// Calc the two upper point
			x1, y1 := getCirclePoint(startAngle, width/2, height/2, width/2)
			x2, y2 := getCirclePoint(endAngle, width/2, height/2, width/2)

			// Create the part corresponding
			svg += fmt.Sprintf(`<path stroke="black" d="M%d,%d L%d,%d A%d,%d 0 %d,1 %d,%d Z" fill="%s" />`,
				width/2, height/2, x1, y1, width/2, height/2, largeArcFlag, x2, y2, getColor(label))

			// Update for the next part
			startAngle = endAngle
		}
	}

	// Finish and return the SVG
	svg += "</svg>"
	return svg
}

func generateLineChart(data []StringInt, width, height int) string {
	// Get the max and the min value from the list
	var max float32
	var min float32
	for _, v := range data {
		if v.Value > max { max = v.Value }
		if v.Value < min { min = v.Value }
	}
	range_value := (max - min) / 5

	// Create the main SVG tag
	svg := fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`, 
		width + 5, height, width, height)

	// Add the line to get the values
	for i := 0; i <= 5; i++ {
		val := range_value * float32(5 - i) + min
		y := float64(height - 15) - float64((val - min) * float32(height - 15) / (max - min)) + 15

		svg += fmt.Sprintf(`<line x1="0" y1="%f" x2="%d" y2="%f" stroke="gray" />`, y, width, y)
		svg += fmt.Sprintf(`<text x="%d" y="%f" font-size="14">%d</text>`, width-35, y-2, int(val))
	}

	// Draw the line
	svg += `<polyline class="point" points="`
	for i, v := range data {
		x := width / Max_RAM_Value * (i + 1) - 10
		y := float32(height) - ((v.Value - min) * float32(height) / (max - min))

		svg += fmt.Sprintf(`%d,%d `, x, int(y))
	}
	svg += `" stroke="blue" fill="none" />`

	// Finish and return the SVG
	svg += `<marker id="circle" markerWidth="12" markerHeight="12" refX="6" refY="6" markerUnits="userSpaceOnUse"> <circle cx="6" cy="6" r="3" stroke-width="2" stroke="context-stroke" fill="context-fill"  /> </marker>`
	svg += "</svg>"
	return svg
}


func getCirclePoint(angle float64, cx, cy, radius int) (int, int) {
	rad := angle * 3.141592653589793 / 180			// Convert the angle in radians
	// Calculate the coordinates with some math i don't understand
	x := cx + int(float64(radius)*math.Cos(rad))
	y := cy + int(float64(radius)*math.Sin(rad))
	return x, y
}

func getColor(name string) string {
	switch name {
		case "Archived": return "#97ae12"
		case "Finished": return "#1232ce"
		case "Fonctionnal": return "#06b41b"
		case "Started": return "#d41212"
		case "Idea": return "#898784"
	}

	return "#000"
}
