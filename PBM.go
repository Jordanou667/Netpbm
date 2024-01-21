package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPPM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	// Open the file for reading
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err // Return an error if file opening fails
	}
	defer file.Close()

	// Create a scanner to read from the file
	scanner := bufio.NewScanner(file)

	// Read the first line to determine the PBM format
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P1" && line != "P4" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	magicNumber := line

	// Read and skip comment lines
	for scanner.Scan() {
        if scanner.Text()[0] == '#' {
            continue
        }
        break

    }

	// Read dimensions (width and height)
    dimension = scanner.Text()
    res := strings.Split(dimension, " ")
    height, _  := strconv.Atoi(res[0])
    width, _  := strconv.Atoi(res[1])
	
	// Read binary data
	var pbm *PBM

	// Check if the PBM image format is "P1"
	if magicNumber == "P1" {
		// Read pixel data for P1 format
		data := make([][]bool, height)
		for i := range data {
			data[i] = make([]bool, width)
		}

		// Iterate over each line of the PBM file to read pixel data
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			hori := strings.Fields(line)
			// Iterate over each field to extract binary pixel values
			for j := 0; j < width; j++ {
				verti, _ := strconv.Atoi(hori[j])
				if verti == 1 {
					data[i][j] = true
				}
			}
		}
		
		// Create a new PBM structure with the read data
		pbm = &PBM{
			data:        data,
			width:       width,
			height:      height,
			magicNumber: magicNumber,
		}
		fmt.Printf("%+v\n", PBM{data, width, height, magicNumber})
	}
	return pbm, nil
}

// Size returns the width and height of the PBM image.
func (pbm *PBM) Size() (int, int) {
	// The function returns the width and height of the PBM image.
    // Return the width and height of the PBM image
	return pbm.width, pbm.height
}

// At retrieves the binary value of a pixel at the specified coordinates in the PBM image.
func (pbm *PBM) At(x, y int) bool {
	 // The function retrieves the binary value (true or false) of a pixel at the specified coordinates.
    // Return the binary value of the pixel at the given coordinates
	return pbm.data[y][x]
}

// Set sets the value of a pixel at the specified coordinates in the PBM image.
func (pbm *PBM) Set(x, y int, value bool) {
	// The function sets the binary value (true or false) of a pixel at the specified coordinates.
    // Update the binary value of the pixel at the given coordinates
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file with the specified filename
func (pbm *PBM) Save(filename string) error {
	// Open the file for writing
    fileSave, error := os.Create(filename)
    if error != nil {
        return error // Return an error if file creation fails
    }
	// Ensure the file is closed when the function exits
	defer fileSave.Close()
	
	// Write the PBM header to the file
    fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	 // Check if the image format is P1
	if pbm.magicNumber == "P1" {
		// Write the pixel data to the file
		for _, i := range pbm.data {
			for _, j := range i {
				if j {
					fmt.Fprint(fileSave, "1 ")
				} else {
					fmt.Fprint(fileSave, "0 ")
				}
			}
			// Move to a new line after each row of pixels
			fmt.Fprintln(fileSave)
		}
	}
	// Return nil to indicate successful save
    return nil
}

// Invert inverts the values of the pixels in the PBM image.
func (pbm *PBM) Invert(){
	// Iterate through each pixel in the PBM image
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			// Check the current value of the pixel
			if pbm.data[i][j] {
				// If true, set it to false
				pbm.data[i][j] = false
			}else if !pbm.data[i][j] {
				// If false, set it to true
				pbm.data[i][j] = true
			}
		}
	}
}

// Flip vertically flips the PBM image.
func (pbm *PBM) Flip() {
	// The function performs a vertical flip by swapping the pixel columns from top to bottom.
    // Iterate through each row of the PBM data and swap corresponding columns from top to bottom
    for _, height := range pbm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}

// Flop horizontally flips the PBM image.
func (pbm *PBM) Flop(){
	// The function performs a horizontal flip by swapping the pixel rows from left to right.
    // Iterate through the first half of the rows, swapping with their corresponding rows from the end
    for i, j := 0, len(pbm.data)-1; i < j; i, j = i+1, j-1 {
        pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
    }
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string){
	// This function allows external modification of the magic number of the PBM image.
	pbm.magicNumber = magicNumber
}
