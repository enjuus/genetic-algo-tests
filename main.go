package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type TwoD struct {
	DNA     *image.RGBA
	Fitness int64
}

var PoolSize = 30
var PopSize = 250
var FitnessLimit int64 = 7500
var MutationRate = 0.0003

func main() {
	start := time.Now()
	rand.Seed(time.Now().UTC().UnixNano())
	target := load("./yui.png")
	population := createPopulation(target)

	found := false
	generation := 0
	for !found {
		generation++
		bestWaifu := getBest(population)
		if bestWaifu.Fitness < FitnessLimit {
			found = true
		} else {
			pool := createPool(population, target)
			population = naturalSelection(pool, population, target)
			if generation%100 == 0 {
				sofar := time.Since(start)
				fmt.Printf("\nTime taken so far: %s | generation: %d | fitness: %d | pool size: %d", sofar, generation, bestWaifu.Fitness, len(pool))
				save("./evolved"+strconv.Itoa(generation)+".png", bestWaifu.DNA)
				fmt.Println()
			}
		}
	}
	elapsed := time.Since(start)
	fmt.Println("\nTotal time taken: %s\n", elapsed)

}
func createBestWaifu(target *image.RGBA) (individual TwoD) {
	individual = TwoD{
		DNA:     createRandomImageFrom(target),
		Fitness: 0,
	}
	individual.calcFitness(target)
	return
}

func save(filePath string, rgba *image.RGBA) {
	imgFile, err := os.Create(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot create file: ", err)
	}

	png.Encode(imgFile, rgba.SubImage(rgba.Rect))
}

func load(filePath string) *image.RGBA {
	imgFile, err := os.Open(filePath)
	defer imgFile.Close()
	if err != nil {
		fmt.Println("Cannot read file: ", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Cannot decode file: ", err)
	}
	return img.(*image.RGBA)
}

func createPopulation(target *image.RGBA) (population []TwoD) {
	population = make([]TwoD, PopSize)
	for i := 0; i < PopSize; i++ {
		population[i] = createBestWaifu(target)
	}
	return
}

func getBest(population []TwoD) TwoD {
	best := int64(0)
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}

func naturalSelection(pool []TwoD, population []TwoD, target *image.RGBA) []TwoD {
	next := make([]TwoD, len(population))

	for i := 0; i < len(population); i++ {
		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		b := pool[r2]

		child := breed(a, b)
		child.mutate()
		child.calcFitness(target)

		next[i] = child
	}
	return next
}

func createRandomImageFrom(img *image.RGBA) (created *image.RGBA) {
	pix := make([]uint8, len(img.Pix))
	rand.Read(pix)
	created = &image.RGBA{
		Pix:    pix,
		Stride: img.Stride,
		Rect:   img.Rect,
	}
	return created
}

func (t *TwoD) calcFitness(target *image.RGBA) {
	difference := diff(t.DNA, target)
	if difference == 0 {
		t.Fitness = 1
	}
	t.Fitness = difference
}

func diff(a, b *image.RGBA) (d int64) {
	d = 0
	for i := 0; i < len(a.Pix); i++ {
		d += int64(sqareDifference(a.Pix[i], b.Pix[i]))
	}
	return int64(math.Sqrt(float64(d)))
}

func sqareDifference(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}

func createPool(population []TwoD, target *image.RGBA) (pool []TwoD) {
	pool = make([]TwoD, 0)
	// get top fitting waifu
	sort.SliceStable(population, func(i, j int) bool {
		return population[i].Fitness < population[j].Fitness
	})
	top := population[0 : PoolSize+1]
	// if there is no difference between top organisms the population is stable
	// and we cant generate a new pool to create the next generation
	if top[len(top)-1].Fitness-top[0].Fitness == 0 {
		pool = population
		return
	}
	for i := 0; i < len(top)-1; i++ {
		num := top[PoolSize].Fitness - top[i].Fitness
		for n := int64(0); n < num; n++ {
			pool = append(pool, top[i])
		}
	}
	return
}

func breed(d1 TwoD, d2 TwoD) TwoD {
	pix := make([]uint8, len(d1.DNA.Pix))
	child := TwoD{
		DNA: &image.RGBA{
			Pix:    pix,
			Stride: d1.DNA.Stride,
			Rect:   d1.DNA.Rect,
		},
		Fitness: 0,
	}
	mid := rand.Intn(len(d1.DNA.Pix))
	for i := 0; i < len(d1.DNA.Pix); i++ {
		if i > mid {
			child.DNA.Pix[i] = d1.DNA.Pix[i]
		} else {
			child.DNA.Pix[i] = d2.DNA.Pix[i]
		}
	}
	return child
}

func (t *TwoD) mutate() {
	for i := 0; i < len(t.DNA.Pix); i++ {
		if rand.Float64() < MutationRate {
			t.DNA.Pix[i] = uint8(rand.Intn(255))
		}
	}
}
