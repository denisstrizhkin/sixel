package sixel

import (
	"log"
	"math/rand"
	"time"
)

type Pixel struct {
	R       uint32
	G       uint32
	B       uint32
	A       uint32
	W       uint32
	Cluster int
}

func (p *Pixel) Dist(q Pixel) uint32 {
	dx := q.R - p.R
	dy := q.G - p.G
	dz := q.B - p.B
	return dx*dx + dy*dy + dz*dz
}

func DistArgMin(cntrs []Pixel, c Pixel) int {
	min := cntrs[0].Dist(c)
	min_i := 0
	for i := 1; i < len(cntrs); i++ {
		d := cntrs[i].Dist(c)
		if d < min {
			min = d
			min_i = i
		}
	}
	return min_i
}

func CompareCntrs(a []Pixel, b []Pixel) bool {
	for i := range a {
		if a[i].R != b[i].R || a[i].G != b[i].G || a[i].B != b[i].B {
			return false
		}
	}
	return true
}

func WeightPixels(pixels []Pixel) (map[Pixel]uint32, []Pixel) {
	wmap := make(map[Pixel]uint32)
	for _, p := range pixels {
		wmap[p]++
	}
	weighted := make([]Pixel, 0, len(wmap))
	for p := range wmap {
		weighted = append(weighted, p)
	}
	log.Println(len(pixels), len(weighted))
	return wmap, weighted
}

func Clusterize(pixels []Pixel, k int, epochs int) ([]Pixel, map[Pixel]int) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	wmap, weighted := WeightPixels(pixels)
	clusterMap := make(map[Pixel]int)

	new_cntrs := make([]Pixel, k)
	cntrs := make([]Pixel, 0, k)
	for range k {
		cntrs = append(cntrs, weighted[rnd.Intn(len(weighted))])
	}

	for range epochs {
		for i := range new_cntrs {
			new_cntrs[i] = Pixel{}
		}

		for pixel, W := range wmap {
			min_i := DistArgMin(cntrs, pixel)
			clusterMap[pixel] = min_i
			new_cntrs[min_i].R += pixel.R * W
			new_cntrs[min_i].G += pixel.G * W
			new_cntrs[min_i].B += pixel.B * W
			new_cntrs[min_i].W += W
		}

		for i := range new_cntrs {
			if new_cntrs[i].W == 0 {
				continue
			}
			new_cntrs[i].R /= new_cntrs[i].W
			new_cntrs[i].G /= new_cntrs[i].W
			new_cntrs[i].B /= new_cntrs[i].W
		}

		if CompareCntrs(cntrs, new_cntrs) {
			log.Println("converged")
			return cntrs, clusterMap
		}
		copy(cntrs, new_cntrs)
	}

	return cntrs, clusterMap
}
