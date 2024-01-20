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
	file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    scanner.Scan()
    line := scanner.Text()
    line = strings.TrimSpace(line)
    if line != "P3" && line != "P6" {
        return nil, fmt.Errorf("Not a Portable Bitmap file: bad magic number %s", line)
    }
    magicNumber := line

	// Read dimensions
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			continue
		}
		break
	}

	dimension := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(dimension[0])
	height, _ := strconv.Atoi(dimension[1])

	scanner.Scan()
    maxValue, _ := strconv.Atoi(scanner.Text())
	// Read pixel data
	var ppm *PPM
	if magicNumber == "P3" {
		data := make([][]Pixel, height)
		for i := range data {
			data[i] = make([]Pixel, width)
		}

		
			// Read pixel data for P3 format
		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			values := strings.Fields(line)
			for j := 0; j < width; j++ {
				if j == 0 {
					r, _ := strconv.Atoi(values[0])
					g, _ := strconv.Atoi(values[1])
					b, _ := strconv.Atoi(values[2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				} else {
					r, _ := strconv.Atoi(values[3*j])
					g, _ := strconv.Atoi(values[3*j+1])
					b, _ := strconv.Atoi(values[3*j+2])
					data[i][j] = Pixel{uint8(r), uint8(g), uint8(b)}
				}
			}
		}

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

func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

func (ppm *PPM) At(x, y int) Pixel{
	return ppm.data[y][x]
}

func (ppm *PPM) Set(x, y int, value Pixel){
	ppm.data[y][x] = value
}

func (ppm *PPM) Save(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Écrire l'en-tête PPM dans le fichier
    _, err = fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
    if err != nil {
        return err
    }

    // Écrire les données des pixels dans le fichier
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            pixel := ppm.data[i][j]
            _, err := fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
            if err != nil {
                return err
            }
        }
        _, err := fmt.Fprintln(file) // Nouvelle ligne après chaque ligne de pixels
        if err != nil {
            return err
        }
    }
    
    return nil
}

func (ppm *PPM) Flip() {
    for _, height := range ppm.data {
        for i, j := 0, len(height)-1; i < j; i, j = i+1, j-1 {
            height[i], height[j] = height[j], height[i]
        }
    }
}

func (ppm *PPM) Invert() {
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            pixel := ppm.data[i][j]
            invertedPixel := Pixel{
            R: ppm.max - pixel.R,
            G: ppm.max - pixel.G,
            B: ppm.max - pixel.B,
            }
            ppm.data[i][j] = invertedPixel
        }
    }
}

func (ppm *PPM) Flop(){
    for i, j := 0, len(ppm.data)-1; i < j; i, j = i+1, j-1 {
        ppm.data[i], ppm.data[j] = ppm.data[j], ppm.data[i]
    }
}

func (ppm *PPM) SetMagicNumber(magicNumber string) {
    ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
    // Vérifier si la nouvelle valeur maximale est différente de la valeur actuelle
    if maxValue == ppm.max {
        return // Pas besoin de faire des modifications si la valeur maximale est la même
    }

    // Calculer le facteur de proportion pour ajuster les valeurs des pixels
    scaleFactor := float64(maxValue) / float64(ppm.max)

    // Ajuster les données des pixels en fonction de la nouvelle valeur maximale
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            pixel := ppm.data[i][j]
            adjustedPixel := Pixel{
                R: uint8(float64(pixel.R) * scaleFactor),
                G: uint8(float64(pixel.G) * scaleFactor),
                B: uint8(float64(pixel.B) * scaleFactor),
            }
            ppm.data[i][j] = adjustedPixel
        }
    }

    // Mettre à jour la nouvelle valeur maximale
    ppm.max = maxValue
}

func (ppm *PPM) Rotate90CW(){
    rotate := make([][]Pixel, ppm.width)
	for i := range rotate {
		rotate[i] = make([]Pixel, ppm.height)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotate[j][ppm.height-1-i] = ppm.data[i][j]
		}
	}

	ppm.data = rotate
	ppm.height, ppm.width = ppm.width, ppm.height
}

func (ppm *PPM) ToPGM() *PGM {
    // Créer une nouvelle matrice pour les données PGM
    pgmData := make([][]uint8, ppm.height)
    for i := range pgmData {
        pgmData[i] = make([]uint8, ppm.width)
    }

    // Convertir les pixels PPM en niveaux de gris PGM
    for i := 0; i < ppm.height; i++ {
        for j := 0; j < ppm.width; j++ {
            pixel := ppm.data[i][j]
            averageValue := uint8((uint32(pixel.R) + uint32(pixel.G) + uint32(pixel.B)) / 3)
            pgmData[i][j] = averageValue
        }
    }

    // Créer une nouvelle instance de la structure PGM
    pgm := &PGM{
        data:        pgmData,
        width:       ppm.width,
        height:      ppm.height,
        magicNumber: "P2",
        max:         ppm.max,
    }

    return pgm
}

func (ppm *PPM) ToPBM() *PBM {
    // Convertir l'image PPM en PGM
    pgm := ppm.ToPGM()

    // Convertir l'image PGM en PBM
    var data [][]bool
    data = make([][]bool, pgm.height)
    for i := range data {
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

    // Créer une nouvelle instance de la structure PBM
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

func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

    for x := 0; x < ppm.height; x++ {
        for y := 0; y < ppm.width; y++ {
            dx := float64(x) - float64(center.X)
            dy := float64(y) - float64(center.Y)
            distance := math.Sqrt(dx*dx + dy*dy)

            if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
                ppm.Set(x, y, color)
            }
        }
    }
    ppm.Set(center.X-(radius-1), center.Y, color)
    ppm.Set(center.X+(radius-1), center.Y, color)
    ppm.Set(center.X, center.Y+(radius-1), color)
    ppm.Set(center.X, center.Y-(radius-1), color)
}

func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
    ppm.DrawCircle(center, radius, color)

    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        for j := 0; j < ppm.width; j++ {
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        if number_points > 1 {
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}

func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
    // Draw the three sides of the triangle
    ppm.DrawLine(p1, p2, color)
    ppm.DrawLine(p2, p3, color)
    ppm.DrawLine(p3, p1, color)
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
    ppm.DrawTriangle(p1, p2, p3, color)

    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        for j := 0; j < ppm.width; j++ {
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        if number_points > 1 {
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
    p2 := Point{p1.X + width, p1.Y}
    p3 := Point{p1.X + width, p1.Y + height}
    p4 := Point{p1.X, p1.Y + height}

    if p1.X < 0 {
        width += p1.X
        p1.X = 0
    }
    if p1.Y < 0 {
        height += p1.Y
        p1.Y = 0
    }
    if p1.X+width > ppm.width {
        width = ppm.width - p1.X
        ppm.DrawLine(p4, p1, color)
    } else {
        ppm.DrawLine(p4, p1, color)
        ppm.DrawLine(p2, p3, color)
    }
    if p1.Y+height > ppm.height {
        height = ppm.height - p1.Y
        ppm.DrawLine(p1, p2, color)
    } else {
        ppm.DrawLine(p1, p2, color)
        ppm.DrawLine(p3, p4, color)
    }

    // Vérifier si les dimensions du rectangle sont maintenant valides
    if width <= 0 || height <= 0 {
        // Dimensions invalides, ne rien faire
        return
    }
}

func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
    ppm.DrawRectangle(p1, width, height, color)

    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        for j := 0; j < ppm.width; j++ {
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        if number_points > 1 {
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

func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
    ppm.DrawPolygon(points, color)
    for i := 0; i < ppm.height; i++ {
        var positions []int
        var number_points int
        for j := 0; j < ppm.width; j++ {
            if ppm.data[i][j] == color {
                number_points += 1
                positions = append(positions, j)
            }
        }
        if number_points > 1 {
            for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
                ppm.data[i][k] = color

            }
        }
    }
}