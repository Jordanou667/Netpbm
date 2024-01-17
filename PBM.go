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

func ReadPBM(filename string) (*PBM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lecture de la première ligne pour obtenir le magic number
	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P1" && line != "P4" {
		return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
	}
	magicNumber := line

	// Lecture des dimensions
	for scanner.Scan() {
        if scanner.Text()[0] == '#' {
            continue
        }
        break

    }

    dimension = scanner.Text()
    res := strings.Split(dimension, " ")
    height, _  := strconv.Atoi(res[0])
    width, _  := strconv.Atoi(res[1])
	
	// Lecture des données binaires
	var pbm *PBM

	if magicNumber == "P1" {
		data := make([][]bool, height)
		for i := range data {
			data[i] = make([]bool, width)
		}

		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			hori := strings.Fields(line)
			for j := 0; j < width; j++ {
				verti, _ := strconv.Atoi(hori[j])
				if verti == 1 {
					data[i][j] = true
				}
			}
		}
		
		pbm = &PBM{
			data:        data,
			width:       width,
			height:      height,
			magicNumber: magicNumber,
		}
		fmt.Printf("%+v\n", PBM{data, width, height, magicNumber})

	// } else if magicNumber == "P4" {
	// 	// Lire le format P4 (binaire)
	// 	expectedBytesPerRow := (width + 7) / 8
	// 	data := make([][]bool, height)
	// 	for i := range data {
	// 		data[i] = make([]bool, width)
	// 	}
		
	// 	for y := 0; y < height; y++ {
	// 		row := make([]byte, expectedBytesPerRow)
	// 		n, err := file.Read(row)
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				return nil, fmt.Errorf("unexpected end of file at line %d", y)
	// 			}
	// 			return nil, fmt.Errorf("error reading pixel data at line %d: %v", y, err)
	// 		}
	// 		if n < expectedBytesPerRow {
	// 			return nil, fmt.Errorf("unexpected end of file at line %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
	// 		}
		
	// 		for x := 0; x < width; x++ {
	// 			byteIndex := x / 8
	// 			bitIndex := 7 - (x % 8)
	
	// 			// Extract the bit from the byte
	// 			bitValue := (int(row[byteIndex]) >> bitIndex) & 1
	
	// 			data[y][x] = bitValue != 0
	// 		}
	// 	}
	
	// 	pbm = &PBM{
	// 		data:        data,
	// 		width:       width,
	// 		height:      height,
	// 		magicNumber: magicNumber,
	// 	}
	// 	fmt.Printf("%+v\n", *pbm)
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
    fileSave, error := os.Create(filename)
    if error != nil {
        return error
    }
	defer fileSave.Close()
	
    fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	if pbm.magicNumber == "P1" {
		for _, i := range pbm.data {
			for _, j := range i {
				if j {
					fmt.Fprint(fileSave, "1 ")
				} else {
					fmt.Fprint(fileSave, "0 ")
				}
			}
			fmt.Fprintln(fileSave)
		}
	// } else if pbm.magicNumber == "P4" {
	// 	// Écrire le format P4 (binaire)
	// 	for y := 0; y < pbm.height; y++ {
	// 		var currentByte byte
	// 		for x := 0; x < pbm.width; x++ {
	// 			bitIndex := 7 - (x % 8)
	// 			bitValue := 0
	// 			if pbm.data[y][x] {
	// 				bitValue = 1
	// 			}
	// 			// Mettre à jour le bit approprié dans l'octet
	// 			currentByte |= byte(bitValue << bitIndex)
	
	// 			if (x+1)%8 == 0 || x == pbm.width-1 {
	// 				_, err := fileSave.Write([]byte{currentByte})
	// 				if err != nil {
	// 					return fmt.Errorf("erreur d'écriture des données binaires à la ligne %d : %v", y, err)
	// 				}
	// 				currentByte = 0
	// 			}
	// 		}
	// 	}
	}	
    return nil
}

func (pbm *PBM) Flip() {
    for _, height := range pbm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}

func (pbm *PBM) Invert(){
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				pbm.data[i][j] = false
			}else if !pbm.data[i][j] {
				pbm.data[i][j] = true
			}
		}
	}
}


func (pbm *PBM) Flop(){
    for i, j := 0, len(pbm.data)-1; i < j; i, j = i+1, j-1 {
        pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
    }
}

func (pbm *PBM) SetMagicNumber(magicNumber string){
	pbm.magicNumber = magicNumber
}

// func main() {
//     pbm, _ := ReadPBM("testImages/pbm/testP1.pbm")
//     // (*PBM).Size(&PBM{})
//     pbm.Save("testImages/pbm/save.pbm")
// 	// fmt.Println("\n")

// 	// pbm.SetMagicNumber("P4")
// 	pbm.Flip()
// 	fmt.Println("Flip:", pbm.data)
// 	// fmt.Println("\n")

// 	pbm.Flop()
// 	fmt.Println("Flop:", pbm.data)
// }