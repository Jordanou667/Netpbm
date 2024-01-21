package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data        [][]uint8
	width, height int
	magicNumber  string
	max          uint8
}

// ReadPPM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
    // Open the file for reading
    file, err := os.Open(filename)
    if err != nil {
        return nil, err // Return an error if file opening fails
    }
    defer file.Close()

    // Create a scanner to read from the file
    scanner := bufio.NewScanner(file)

    // Read the first line to determine the PGM format
    scanner.Scan()
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if line != "P2" && line != "P5" {
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

    // Read the max value (not declared in your initial code)
    scanner.Scan()
    maxValue, _ := strconv.Atoi(scanner.Text())

    // Read binary data
    var pgm *PGM

    // Check if the PGM image format is "P2"
    if magicNumber == "P2" {
        // Read pixel data for P2 format
        data := make([][]uint8, height)
        for i := range data {
            data[i] = make([]uint8, width)
        }

        // Iterate over each row in the PGM image
        for i := 0; i < height; i++ {
            scanner.Scan()
            line := scanner.Text()
            hori := strings.Fields(line)
            // Iterate over each pixel value in the row
            for j := 0; j < width; j++ {
                pixel, _ := strconv.Atoi(hori[j])
                data[i][j] = uint8(pixel)
            }
        }

        // Create a new instance of the PGM structure
        pgm = &PGM{
            data:        data,
            width:       width,
            height:      height,
            magicNumber: magicNumber,
            max:         uint8(maxValue),
        }
        fmt.Printf("%+v\n", PGM{data, width, height, magicNumber, uint8(maxValue)})
    }
    return pgm, nil
}

// Size returns the width and height of the PGM image
func (pgm *PGM) Size() (int, int) {
    // The function returns the width and height of the PGM image.
    // Return the width and height of the PGM image
	return pgm.width, pgm.height
}

// At retrieves the intensity value of a pixel at the specified coordinates in the PGM image.
func (pgm *PGM) At(x, y int) uint8{
    // The function retrieves the intensity value of a pixel at the specified coordinates.
    // Return the intensity value of the pixel at the given coordinates
	return pgm.data[y][x]
}

// Set sets the value of a pixel at the specified coordinates in the PGM image.
func (pgm *PGM) Set(x, y int, value uint8){
    // The function sets the binary value (true or false) of a pixel at the specified coordinates.
    // Update the value of the pixel at the given coordinates
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file with the specified filename
func (pgm *PGM) Save(filename string) error {
    // Open the file for writing
    fileSave, error := os.Create(filename)
    if error != nil {
        return error // Return an error if file creation fails
    }
    // Ensure the file is closed when the function exits
    defer fileSave.Close()

    // Write the PGM header to the file
    fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

     // Check if the image format is P2
    if pgm.magicNumber == "P2" {
        // Write the pixel data to the file
        for i := 0; i < pgm.height; i++ {
            for j := 0; j < pgm.width; j++ {
                fmt.Fprintf(fileSave, "%v ", pgm.data[i][j])
            }
            fmt.Fprintf(fileSave, "\n")
        }
        // Move to a new line after each row of pixels
        fmt.Fprintln(fileSave)
    }
    // Return nil to indicate successful save
    return nil
}

// Invert inverts the intensity values of the pixels in the PGM image.
func (pgm *PGM) Invert() {
    // Check if the PGM image dimensions are valid
    if pgm.width == 0 || pgm.height == 0 {
        return // Return if the image dimensions are invalid
    }

    // Iterate through each pixel in the PGM image
    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
            // Invert the intensity value by subtracting it from the maximum intensity
            pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
        }
    }
}

// Flip vertically flips the PGM image.
func (pgm *PGM) Flip() {
    // The function performs a vertical flip by swapping the pixel columns from top to bottom.
    // Iterate through each row of the PGM data and swap corresponding columns from top to bottom
    for _, height := range pgm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}

// Flop horizontally flips the PGM image.
func (pgm *PGM) Flop(){
    // The function performs a horizontal flip by swapping the pixel rows from left to right.
    // Iterate through the first half of the rows, swapping with their corresponding rows from the end
    for i, j := 0, len(pgm.data)-1; i < j; i, j = i+1, j-1 {
        pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
    }
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string){
    // This function allows external modification of the magic number of the PGM image.
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
    // Set the multiplicator
    multiplicator := float64(maxValue) / float64(pgm.max)
    // ppm.max becomes our new max valuea
    pgm.max = uint8(maxValue)

    // Updates pixel values with the new max value
    for i := range pgm.data {
        for j := range pgm.data[i] {
            //Modifies the value of each pixel proportionally
            pgm.data[i][j] = uint8(float64(pgm.data[i][j]) * float64(multiplicator))
        }
    }

}

// Rotate90CW rotates the PGM image 90 degrees clockwise.
func (pgm *PGM) Rotate90CW(){
    // Create a new matrix to store rotated data
    rotate := make([][]uint8, pgm.width)
	for i := range rotate {
		rotate[i] = make([]uint8, pgm.height)
	}

    // Iterate through each pixel of the original PGM image
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
            // Rotate each pixel 90 degrees clockwise and assign it to the new matrix
			rotate[j][pgm.height-1-i] = pgm.data[i][j]
		}
	}

    // Update the PGM image data with the rotated matrix
	pgm.data = rotate
    // Swap height and width values to reflect the rotation
	pgm.height, pgm.width = pgm.width, pgm.height
}

// ToPBM converts a PGM image to a PBM image by thresholding based on intensity.
func (pgm *PGM) ToPBM() *PBM {
    // Create a new matrix for PBM data
    var data [][]bool
    data = make([][]bool, pgm.height)
    for i := 0; i < pgm.width; i++ {
        data[i] = make([]bool, pgm.width)
    }

    // Convert PGM pixels to PBM binary values based on intensity threshold
    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
            // Thresholding: If intensity is greater than half of max, set to true; otherwise, set to false
            if pgm.data[i][j] > pgm.max/2 {
                data[i][j] = true
            } else {
                data[i][j] = false
            }
        }
    }

    // Create a new instance of the PBM structure
    return &PBM{
        data:        data,
        width:       pgm.width,
        height:      pgm.height,
        magicNumber: "P1",
    }
}
