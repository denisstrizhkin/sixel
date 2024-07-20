package sixel

import (
	"math/rand"
	"time"
	"log"
)

type Pixel struct {
	R uint32
	G uint32
	B uint32
	A uint32
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
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Clusterize(pixels []Pixel, k int,  epochs int) []Pixel {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	cntrs := make([]Pixel, 0, k)
	for range k {
		cntrs = append(cntrs, pixels[rnd.Intn(len(pixels))])
	}
	
	var new_cntrs []Pixel
	for range epochs {
		new_cntrs = make([]Pixel, k)
		for i, pixel := range pixels {
			min_i := DistArgMin(cntrs, pixel)
			pixels[i].Cluster = min_i
			new_cntrs[min_i].R += pixel.R
			new_cntrs[min_i].G += pixel.G
			new_cntrs[min_i].B += pixel.B
			new_cntrs[min_i].A++	
		}

		for i := range new_cntrs {
			if new_cntrs[i].A == 0 {
				continue
			}
			new_cntrs[i].R /= new_cntrs[i].A
			new_cntrs[i].G /= new_cntrs[i].A
			new_cntrs[i].B /= new_cntrs[i].A
		}

		if CompareCntrs(cntrs, new_cntrs) {
			log.Println("converged")
			return cntrs;
		}
		cntrs = new_cntrs
	}

	return cntrs;
}
