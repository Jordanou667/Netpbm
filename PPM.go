package Netpbm

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// type PPM struct {
// 	data        [][]Pixel
// 	width, height int
// 	magicNumber  string
// 	max          int
// }

// type Pixel struct {
// 	R, G, B uint8
// }

// // ReadPPM reads a PPM image from a file and returns a struct that represents the image.
// func ReadPPM(filename string) (*PPM, error) {
// 	file, err := os.Open(filename)
//     if err != nil {
//         return nil, err
//     }
//     defer file.Close()

//     scanner := bufio.NewScanner(file)

//     scanner.Scan()
//     line := scanner.Text()
//     line = strings.TrimSpace(line)
//     if line != "P3" && line != "P6" {
//         return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
//     }
//     magicNumber := line

// 	// Read dimensions
// 	for scanner.Scan() {
// 		if strings.HasPrefix(scanner.Text(), "#") {
// 			continue
// 		}
// 		break
// 	}

// 	dimension := strings.Fields(scanner.Text())
// 	width, _ := strconv.Atoi(dimension[0])
// 	height, _ := strconv.Atoi(dimension[1])

// 	scanner.Scan()
//     maxValue, _ := strconv.Atoi(scanner.Text())

// 	// Read pixel data
// 	data := make([][]Pixel, height)
// 	for i := range data {
// 		data[i] = make([]Pixel, width)
// 	}

	
// 		// Read pixel data for P3 format
// 	for i := 0; i < height; i++ {
// 		scanner.Scan()
// 			line := scanner.Text()
// 			values := strings.Fields(line)
// 		for j := 0; j < width; j++ {
// 			r, _ := strconv.Atoi(values[0])
// 			g, _ := strconv.Atoi(values[1])
// 			b, _ := strconv.Atoi(values[2])
// 			data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
// 		}
// 	}

// 	ppm := &PPM{
// 		data:        data,
// 		width:       width,
// 		height:      height,
// 		magicNumber: magicNumber,
// 		max:         maxValue,
// 	}
	
// 	fmt.Printf("%+v\n", PPM{data, width, height, magicNumber, maxValue})
// 	return ppm, nil
// }


// func (ppm *PPM) Size() (int, int) {
// 	return ppm.width, ppm.height
// }

// func (ppm *PPM) At(x, y int) Pixel{
// 	return ppm.data[y][x]
// }

// func (ppm *PPM) Set(x, y int, value Pixel){
// 	ppm.data[y][x] = value
// }

// func (ppm *PPM) Save(filename string) error{
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = fmt.Fprintf(file, "magicNumber: %s\n", ppm.magicNumber)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = fmt.Fprintf(file, "Width: %d\n", ppm.width)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = fmt.Fprintf(file, "Height: %d\n", ppm.height)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = fmt.Fprintf(file, "Max Value: %d\n", ppm.max)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Fprintf(file, "\n")

// 	// for i := 0; i < ppm.height; i++ {
// 	// 	for j := 0; j < ppm.width; j++ {
// 	// 		fmt.Fprintf(file, "%v ", ppm.data[i][j])
// 	// 	}
// 	// 	fmt.Fprintf(file, "\n")
// 	// }
// 	// fmt.Fprintf(file, "\n")

// 	if ppm.magicNumber == "P3" {
// 		// fmt.Fprintf(file, "R, G, B: ")
// 		// Write pixel data for P3 format
// 		for i := 0; i < ppm.height; i++ {
// 			for j := 0; j < ppm.width; j++ {
// 				fmt.Fprintf(file, "R:%d G:%d B:%d ", ppm.data[i][j].R, ppm.data[i][j].G, ppm.data[i][j].B)
// 			}
// 			fmt.Fprintln(file)
// 		}
// 	}

// 	return nil
// }



// func (ppm *PPM) Flop() {
//     for _, height := range ppm.data {
//         for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
//             height[i], height[j] = height[j], height[i]
//         }
//     }
// }

// func (ppm *PPM) Flip(){
//     for i, j := 0, len(ppm.data)-1; i < j; i, j = i+1, j-1 {
//         ppm.data[i], ppm.data[j] = ppm.data[j], ppm.data[i]
//     }
// }

// func (ppm *PPM) SetMagicNumber(magicNumber string) {
//     ppm.magicNumber = magicNumber
// }

// // SetMaxValue sets the max value of the PPM image.
// func (ppm *PPM) SetMaxValue(maxValue uint8) {
//     ppm.max = int(maxValue)

//     // Adjust pixel values if they exceed the new max value
//     for i := 0; i < ppm.height; i++ {
//         for j := 0; j < ppm.width; j++ {
//             if ppm.data[i][j].R > maxValue {
//                 ppm.data[i][j].R = maxValue
//             }
//             if ppm.data[i][j].G > maxValue {
//                 ppm.data[i][j].G = maxValue
//             }
//             if ppm.data[i][j].B > maxValue {
//                 ppm.data[i][j].B = maxValue
//             }
//         }
//     }
// }

// func (ppm *PPM) Rotate90CW(){
//     rotate := make([][]Pixel, ppm.width)
// 	for i := range rotate {
// 		rotate[i] = make([]Pixel, ppm.height)
// 	}

// 	for i := 0; i < ppm.height; i++ {
// 		for j := 0; j < ppm.width; j++ {
// 			rotate[j][ppm.height-1-i] = ppm.data[i][j]
// 		}
// 	}

// 	ppm.data = rotate
// 	ppm.height, ppm.width = ppm.width, ppm.height
// }




// func main() {
//     ppm, _ := ReadPPM("testImages/ppm/testP3.ppm")
//     // (*PBM).Size(&PBM{})
//     ppm.Save("testImages/ppm/save.ppm")
// 	fmt.Println("\n")

// 	// ppm.SetMagicNumber("P5")
// 	ppm.Flip()
// 	fmt.Println("Flip:", ppm.data)
// 	fmt.Println("\n")

// 	ppm.Flop()
// 	fmt.Println("Flop:", ppm.data)
// 	fmt.Println("\n")

// 	ppm.Rotate90CW()
// 	fmt.Println("Rotate:", ppm.data)

// }