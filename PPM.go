package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
    "math"
)

type PPM struct {
	data        [][]Pixel
	width, height int
	magicNumber  string
	max          uint8
}

type Pixel struct {
	R, G, B uint8
}

type Point struct{
    X, Y int
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
    // Open the file for reading
	file, err := os.Open(filename)
    if err != nil {
        return nil, err // Return an error if file opening fails
    }
    defer file.Close()

    // Create a scanner to read from the file
    scanner := bufio.NewScanner(file)

    // Read the first line to determine the PPM format
    scanner.Scan()
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if line != "P3" && line != "P6" {
        return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
    }
    magicNumber := line

	// Read and skip comment lines
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		break
	}

    // Read dimensions (width and height)
	dimension := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimension[0])
	height, _ := strconv.Atoi(dimension[1])

    // Read the maximum pixel value
	scanner.Scan()
    maxValue, _ := strconv.Atoi(scanner.Text())

	// Read pixel data
	var ppm *PPM

    // Check if the PPM image format is "P3"
	if magicNumber == "P3" {
        // Read pixel data for P3 format
		data := make([][]Pixel, height)
		for i := range data {
			data[i] = make([]Pixel, width)
		}

        // Iterate over each row in the PPM image
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			values := strings.Fields(line)
            // Iterate over each pixel value in the row
			for j := 0; j < width; j++ {
				if j == 0 {
                    // Read RGB values for the first pixel in the line
					r, _ := strconv.Atoi(values[0])
					g, _ := strconv.Atoi(values[1])
					b, _ := strconv.Atoi(values[2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				} else {
                    // Read RGB values for subsequent pixels in the line
					r, _ := strconv.Atoi(values[3*j])
					g, _ := strconv.Atoi(values[3*j+1])
					b, _ := strconv.Atoi(values[3*j+2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				}
			}
		}

        // Create a new instance of the PPM structure
		ppm = &PPM{
			data:        data,
			width:       width,
			height:      height,
			magicNumber: magicNumber,
			max:         uint8(maxValue),
		}
		fmt.Printf("%+v\n", PPM{data, width, height, magicNumber, uint8(maxValue)})
	}
	return ppm, nil
}

// Size returns the width and height of the PPM image.
func (ppm *PPM) Size() (int, int) {
    // The function returns the width and height of the PPM image.
    // Return the width and height of the PPM image
	return ppm.width, ppm.height
}

// At retrieves the RGB values of a pixel at the specified coordinates in the PPM image.
func (ppm *PPM) At(x, y int) Pixel{
    // The function retrieves the RGB values of a pixel at the specified coordinates.
    // Return the RGB values of the pixel at the given coordinates
	return ppm.data[y][x]
}

// Set sets the value of a pixel at the specified coordinates in the PPM image.
func (ppm *PPM) Set(x, y int, value Pixel){
    // The function sets the RGB values of a pixel at the specified coordinates.
    // Update the RGB values of the pixel at the given coordinates
	ppm.data[y][x] = value
}

// Save saves the PPM image to a file with the specified filename.
func (ppm *PPM) Save(filename string) error {
    // Open the file for writing
    file, err := os.Create(filename)
    if err != nil {
        return err // Return an error if file creation fails
    }
    // Ensure the file is closed when the function exits
    defer file.Close()

    // Write the PPM header to the file
    _, err = fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
    if err != nil {
        return err
    }

     // Check if the image format is P3
    if ppm.magicNumber == "P3" {
        // Write the pixel data to the file
        for i := 0; i < ppm.height; i++ {
            for j := 0; j < ppm.width; j++ {
                pixel := ppm.data[i][j]
                _, err := fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
                if err != nil {
                    return err
                }
            }
            // Move to a new line after each row of pixels
            _, err := fmt.Fprintln(file)
            if err != nil {
                return err
            }
        }
    }
    // Return nil to indicate successful save
    return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
    // Iterate through each pixel in the PPM image
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            // Get the original pixel value
            pixel := ppm.data[i][j]
            invertedPixel := Pixel{
                // Invert each color component of the pixel by subtracting it from the maximum value
                R: ppm.max - pixel.R,
                G: ppm.max - pixel.G,
                B: ppm.max - pixel.B,
            }
            // Update the pixel with the inverted values
            ppm.data[i][j] = invertedPixel
        }
    }
}

// Flip vertically flips the PPM image.
func (ppm *PPM) Flip() {
    // The function performs a vertical flip by swapping the pixel columns from top to bottom.
    // Iterate through each row of the PPM data and swap corresponding columns from top to bottom
    for _, height := range ppm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}
// Flop horizontally flips the PPM image.
func (ppm *PPM) Flop(){
    // The function performs a horizontal flip by swapping the pixel rows from left to right.
    // Iterate through the first half of the rows, swapping with their corresponding rows from the end
    for i, j := 0, len(ppm.data)-1; i < j; i, j = i+1, j-1 {
        ppm.data[i], ppm.data[j] = ppm.data[j], ppm.data[i]
    }
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
    // This function allows external modification of the magic number of the PPM image.
    ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
    // Check if the new maximum value is different from the current value
    if maxValue == ppm.max {
        return // No need to make modifications if the maximum value is the same
    }

    // Calculate the proportion factor to adjust pixel values
    scaleFactor := float64(maxValue) / float64(ppm.max)

    // Adjust pixel data based on the new maximum value
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            // Get the original pixel value
            pixel := ppm.data[i][j]

            // Adjust each color component of the pixel using the scaleFactor
            adjustedPixel := Pixel{
                R: uint8(float64(pixel.R) * scaleFactor),
                G: uint8(float64(pixel.G) * scaleFactor),
                B: uint8(float64(pixel.B) * scaleFactor),
            }
            // Update the pixel with the adjusted values
            ppm.data[i][j] = adjustedPixel
        }
    }

    // Update the new maximum value
    ppm.max = maxValue
}

// Rotate90CW rotates the PPM image 90 degrees clockwise.
func (ppm *PPM) Rotate90CW(){
    // Create a new matrix to store rotated data
    rotate := make([][]Pixel, ppm.width)
	for i := range rotate {
		rotate[i] = make([]Pixel, ppm.height)
	}

    // Iterate through each pixel of the original PPM image
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
            // Rotate each pixel 90 degrees clockwise and assign it to the new matrix
			rotate[j][ppm.height-1-i] = ppm.data[i][j]
		}
	}

    // Update the PPM image data with the rotated matrix
	ppm.data = rotate
    // Swap height and width values to reflect the rotation
	ppm.height, ppm.width = ppm.width, ppm.height
}

// ToPGM converts the PPM image to a PGM image
func (ppm *PPM) ToPGM() *PGM {
    // Create a new matrix for PGM data
    pgmData := make([][]uint8, ppm.height)
    for i := range pgmData {
        pgmData[i] = make([]uint8, ppm.width)
    }

    // Convert PPM pixels to PGM grayscale
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            // Extract the RGB values of the pixel
            pixel := ppm.data[i][j]
            // Calculate the average grayscale value
            averageValue := uint8((uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)) / 3)
            // Assign the average value to the corresponding pixel in the PGM data matrix
            pgmData[i][j] = averageValue
        }
    }

    // Create a new instance of the PGM structure
    pgm := &PGM{
        data:        pgmData,
        width:       ppm.width,
        height:      ppm.height,
        magicNumber: "P2",
        max:         ppm.max,
    }

    return pgm
}

// ToPBM converts the PPM image to a PBM image.
func (ppm *PPM) ToPBM() *PBM {
    // Convert PPM image to PGM
    pgm := ppm.ToPGM()

    // Convert PGM image to PBM
    var data [][]bool
    data = make([][]bool, pgm.height)
    for i := range data {
        data[i] = make([]bool, pgm.width)
    }

    // Iterate through each pixel in the PGM image
    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
             // Check if the intensity of the pixel is greater than half of the maximum intensity
            // If true, set the corresponding pixel in the PBM image to true (1), else set it to false (0)
            if pgm.data[i][j] > pgm.max/2 {
                data[i][j] = true
            } else {
                data[i][j] = false
            }
        }
    }

    // Create a new instance of the PBM structure
    pbm := &PBM{
        data:        data,
        width:       pgm.width,
        height:      pgm.height,
        magicNumber: "P1",
    }

    return pbm
}

// DrawLine draws a line between two points on the PPM image.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
    // Handle points outside the image bounds
    if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
        // Find the intersection point with the image bounds
        if p1.X < 0 {
            p1.X = 0
        } else if p1.X >= ppm.width {
            p1.X = ppm.width - 1
        }

        if p1.Y < 0 {
            p1.Y = 0
        } else if p1.Y >= ppm.height {
            p1.Y = ppm.height - 1
        }
    }

    dx := p2.X - p1.X
    dy := p2.Y - p1.Y

    // Determine the direction of the line
    var sx, sy int
    if dx > 0 {
        sx = 1
    } else {
        sx = -1
        dx = -dx
    }
    if dy > 0 {
        sy = 1
    }else {
        sy = -1
        dy = -dy
    }

    err := dx - dy

    // Draw the line
    for {
        // Check if the current point is within the image bounds
        if p1.X >= 0 && p1.X < ppm.width && p1.Y >= 0 && p1.Y < ppm.height {
            ppm.data[p1.Y][p1.X] = color
        }

        // Break the loop when the end point is reached
        if p1.X == p2.X && p1.Y == p2.Y {
            break
        }

        e2 := 2 * err
        if e2 > -dy {
            err -= dy
            p1.X += sx
        }
        if e2 < dx {
            err += dx
            p1.Y += sy
        }
    }
}

// DrawCircle draws a tcircle
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

    // Iterate through each pixel in the image.
    for x := 0; x < ppm.height; x++ {
        for y := 0; y < ppm.width; y++ {
            // Calculate the distance from the current pixel to the center of the circle.
            dx := float64(x) - float64(center.X)
            dy := float64(y) - float64(center.Y)
            distance := math.Sqrt(dx*dx + dy*dy)

            // Check that the distance to the center is approximately equal to the specified radius. 
            //The condition "math.Abs(distance-float64(radius)) < 1.0" allows a small margin of error, checking that the distance is less than the specified radius.
            if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
                // If the conditions are met, set the color of the pixel to the specified color.
                ppm.Set(x, y, color)
            }
        }
    }
    // Draw four points around the circle to complete its outline.
    ppm.Set(center.X-(radius-1), center.Y, color)
    ppm.Set(center.X+(radius-1), center.Y, color)
    ppm.Set(center.X, center.Y+(radius-1), color)
    ppm.Set(center.X, center.Y-(radius-1), color)
}

// DrawFilledCircle draws a circle with the specified dimensions and color on the PPM image.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
    //draw the outline of the circle.
    ppm.DrawCircle(center, radius, color)

    // Iterate through each row of the image.
    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        // Iterate through each column of the image.
        for j := 0; j < ppm.width; j++ {
            // Check if the pixel at (i, j) has the specified color.
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        // If there are more than one pixel with the specified color in the current row, fill the gap between the leftmost and rightmost pixels.
        if number_points > 1 {
            // Iterate through the positions to fill the ga
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}

// DrawTriangle draws a triangle
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
    // Draw the three sides of the triangle
    ppm.DrawLine(p1, p2, color)
    ppm.DrawLine(p2, p3, color)
    ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a triangle with the specified dimensions and color on the PPM image.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
    //draw the outline of the triangle.
    ppm.DrawTriangle(p1, p2, p3, color)

     // Iterate through each row of the image.
    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        // Check if the pixel at (i, j) has the specified color.
        for j := 0; j < ppm.width; j++ {
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }

        // If there are more than one pixel with the specified color in the current row, fill the gap between the leftmost and rightmost pixels.
        if number_points > 1 {
            // Iterate through the positions to fill the gap between the leftmost and rightmost pixels.
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}

// DrawRectangle draws a rectangle
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
    // Calculate the coordinates of the other three vertices of the rectangle.
    p2 := Point{p1.X + width, p1.Y}
    p3 := Point{p1.X + width, p1.Y + height}
    p4 := Point{p1.X, p1.Y + height}

    // If the left edge of the rectangle is outside the image, adjust the width and reset the X-coordinate.
    if p1.X < 0 {
        width += p1.X
        p1.X = 0
    }

     // If the top edge of the rectangle is outside the image, adjust the height and reset the Y-coordinate.
    if p1.Y < 0 {
        height += p1.Y
        p1.Y = 0
    }

    // If the right edge of the rectangle is outside the image, adjust the width and draw a line from p4 to p1.
    if p1.X+width > ppm.width {
        width = ppm.width - p1.X
        ppm.DrawLine(p4, p1, color)
    } else {
        // Draw the top and right edges of the rectangle.
        ppm.DrawLine(p4, p1, color)
        ppm.DrawLine(p2, p3, color)
    }

    // If the bottom edge of the rectangle is outside the image, adjust the height and draw a line from p1 to p2.
    if p1.Y+height > ppm.height {
        height = ppm.height - p1.Y
        ppm.DrawLine(p1, p2, color)
    } else {
        // Draw the left and bottom edges of the rectangle.
        ppm.DrawLine(p1, p2, color)
        ppm.DrawLine(p3, p4, color)
    }

    // Vérifier si les dimensions du rectangle sont maintenant valides
    if width <= 0 || height <= 0 {
        // Dimensions invalides, ne rien faire
        return
    }
}

// DrawRectangle draws a rectangle with the specified dimensions and color on the PPM image.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
    // draw the outline of the rectangle
    ppm.DrawRectangle(p1, width, height, color)

    // Iterate through each row of the image.
    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        // Iterate through each column of the image.
        for j := 0; j < ppm.width; j++ {
            // Check if the pixel at (i, j) has the specified color
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        //If there are more than one pixel with the specified color in the current row, fill the gap between the leftmost and rightmost pixels.
        if number_points > 1 {
            // Iterate through the positions to fill the gap between the leftmost and rightmost pixels.
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
        if height > ppm.height && width > ppm.width {
            for k := 0; k < ppm.width; k++ {
                ppm.data[i][k] = color

            }

        }
    }
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
    numPoints := len(points)
    if numPoints < 3 {
        // A polygon must have at least 3 vertices
        return
    }

    // Draw lines between consecutive points to form the polygon
    for i := 0; i < numPoints-1; i++ {
        ppm.DrawLine(points[i], points[i+1], color)
    }

    // Draw the last line connecting the last and first points to close the polygon
    ppm.DrawLine(points[numPoints-1], points[0], color)
}

// DrawFilledPolygon fills the specified polygon with the given color on the PPM image.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
    //draw the outline of the polygon
    ppm.DrawPolygon(points, color)

    // Iterate through each row of the image.
    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        for j := 0; j < ppm.width; j++ {
            // Check if the pixel at (i, j) has the specified color.
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        //fills the space between the leftmost and rightmost pixels if there is more than one pixel with the color specified in the current line
        if number_points > 1 {
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}
