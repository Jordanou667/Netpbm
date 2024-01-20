package Netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
    // "io"
)

type PGM struct {
	data        [][]uint8
	width, height int
	magicNumber  string
	max          uint8
}

func ReadPGM(filename string) (*PGM, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    scanner.Scan()
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if line != "P2" && line != "P5" {
        return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
    }
    magicNumber := line

    // Lecture des dimensions
    for scanner.Scan() {
        if strings.HasPrefix(scanner.Text(), "#") {
            continue
        }
        break
    }

    dimension := strings.Fields(scanner.Text())
    width, _ := strconv.Atoi(dimension[0])
    height, _ := strconv.Atoi(dimension[1])

    // Lecture du max (non déclaré dans votre code initial)
    scanner.Scan()
    maxValue, _ := strconv.Atoi(scanner.Text())

    // Lecture des données binaires
    var pgm *PGM

    if magicNumber == "P2" {
        data := make([][]uint8, height)
        for i := range data {
            data[i] = make([]uint8, width)
        }

        for i := 0; i < height; i++ {
            scanner.Scan()
            line := scanner.Text()
            hori := strings.Fields(line)
            for j := 0; j < width; j++ {
                pixel, _ := strconv.Atoi(hori[j])
                data[i][j] = uint8(pixel)
            }
        }

        pgm = &PGM{
            data:        data,
            width:       width,
            height:      height,
            magicNumber: magicNumber,
            max:         uint8(maxValue),
        }
        fmt.Printf("%+v\n", PGM{data, width, height, magicNumber, uint8(maxValue)})
    
    // if magicNumber == "P5" {
    //     // Read the format P5 (binary)
    //     data := make([][]uint8, height)
    //     for i := range data {
    //         data[i] = make([]uint8, width)
    //     }
    
    //     for y := 0; y < height; y++ {
    //         row := make([]byte, width)
    //         n, err := file.Read(row)
    //         if err != nil {
    //             if err == io.EOF {
    //                 return nil, fmt.Errorf("unexpected end of file at line %d", y)
    //             }
    //             return nil, fmt.Errorf("error reading pixel data at line %d: %v", y, err)
    //         }
    //         if n < width {
    //             return nil, fmt.Errorf("unexpected end of file at line %d, expected %d bytes, got %d", y, width, n)
    //         }
    
    //         for x := 0; x < width; x++ {
    //             data[y][x] = uint8(row[x])
    //         }
    //     }
    
    //     pgm = &PGM{
    //         data:        data,
    //         width:       width,
    //         height:      height,
    //         magicNumber: magicNumber,
    //         max:         maxValue,
    //     }
    //     fmt.Printf("%+v\n", *pgm)
    }
    return pgm, nil
}


func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8{
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8){
	pgm.data[y][x] = value
}


func (pgm *PGM) Save(filename string) error {
    fileSave, error := os.Create(filename)
    if error != nil {
        return error
    }
    defer fileSave.Close()

    fmt.Fprintf(fileSave, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

    if pgm.magicNumber == "P2" {
        for i := 0; i < pgm.height; i++ {
            for j := 0; j < pgm.width; j++ {
                fmt.Fprintf(fileSave, "%v ", pgm.data[i][j])
            }
            fmt.Fprintf(fileSave, "\n")
        }
        fmt.Fprintln(fileSave)

    // if pgm.magicNumber == "P5" {
    //     // Save binary data for P5 format
    //     for y := 0; y < pgm.height; y++ {
    //         for x := 0; x < pgm.width; x++ {
    //             fmt.Fprintf(fileSave, "%c", pgm.data[y][x])
    //         }
    //     }
    }
    return nil
}

func (pgm *PGM) Invert() {
    if pgm.width == 0 || pgm.height == 0 {
        return
    }

    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
            pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
        }
    }
}

func (pgm *PGM) Flip() {
    for _, height := range pgm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}

func (pgm *PGM) Flop(){
    for i, j := 0, len(pgm.data)-1; i < j; i, j = i+1, j-1 {
        pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
    }
}

func (pgm *PGM) SetMagicNumber(magicNumber string){
	pgm.magicNumber = magicNumber
}

func (pgm *PGM) SetMaxValue(maxValue uint8){
	pgm.max = uint8(maxValue)
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(math.Round(float64(pgm.data[i][j]) / float64(pgm.max) * 255))
		}
	}
}

func (pgm *PGM) Rotate90CW(){
    rotate := make([][]uint8, pgm.width)
	for i := range rotate {
		rotate[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotate[j][pgm.height-1-i] = pgm.data[i][j]
		}
	}

	pgm.data = rotate
	pgm.height, pgm.width = pgm.width, pgm.height
}

func (pgm *PGM) ToPBM() *PBM {
    var data [][]bool
    data = make([][]bool, pgm.height)
    for i := 0; i < pgm.width; i++ {
        data[i] = make([]bool, pgm.width)
    }

    for i := 0; i < pgm.height; i++ {
        for j := 0; j < pgm.width; j++ {
            if pgm.data[i][j] > pgm.max/2 {
                data[i][j] = true
            } else {
                data[i][j] = false
            }
        }
    }

    return &PBM{
        data:        data,
        width:       pgm.width,
        height:      pgm.height,
        magicNumber: "P1",
    }
}

// func main() {
//     pgm, _ := ReadPGM("testImages/pgm/testP5.pgm")
//     // (*PBM).Size(&PBM{})
//     pgm.Save("testImagees/pgm/testp5a.pgm")
// 	// fmt.Println("\n")

// 	// pgm.SetMagicNumber("P5")
// 	pgm.Flip()
// 	fmt.Println("Flip:", pgm.data)
// 	// fmt.Println("\n")

// 	pgm.Flop()
// 	fmt.Println("Flop:", pgm.data)
// 	// fmt.Println("\n")

// 	pgm.Rotate90CW()
// 	fmt.Println("Rotate:", pgm.data)
// }