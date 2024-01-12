package main

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
	max          int
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
    max, _ := strconv.Atoi(scanner.Text())

    // Lecture des données binaires
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

    pgm := &PGM{
        data:        data,
        width:       width,
        height:      height,
        magicNumber: magicNumber,
        max:         max,
    }

    fmt.Printf("%+v\n", PGM{data, width, height, magicNumber, max})

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
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "magicNumber: %s\n", pgm.magicNumber)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "Width: %d\n", pgm.width)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "Height: %d\n", pgm.height)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(file, "Max: %d\n", pgm.max)
	if err != nil {
		return err
	}

	fmt.Fprintf(file, "\n")

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Fprintf(file, "%v ", pgm.data[i][j])
		}
		fmt.Fprintf(file, "\n")
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

func (pgm *PGM) SetMaxValue(maxValue uint8){
	pgm.max = int(maxValue)
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++{
			if pgm.data[i][j] > maxValue {
				pgm.data[i][j] = maxValue
			}
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



// func (pgm *PGM) SetMagicNumber(magicNumber string){
// 	fmt.Println(magicNumber)
// }




func main() {
    pgm, _ := ReadPGM("Dos_pgm/test.pgm")
    // (*PBM).Size(&PBM{})
    pgm.Save("Dos_pgm/save.pgm")
	fmt.Println("\n")

	// pgm.SetMagicNumber("P5")
	pgm.Flip()
	fmt.Println("Flip:", pgm.data)
	fmt.Println("\n")

	pgm.Flop()
	fmt.Println("Flop:", pgm.data)
	fmt.Println("\n")

	pgm.Rotate90CW()
	fmt.Println("Rotate:", pgm.data)
}