package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"
)

type XYZ struct {
	X int
	Y int
	Z int
	C int
}

func RGBAToXYZ(c color.RGBA) XYZ {
	return XYZ{
		X: int(c.A),
		Y: int(c.G),
		Z: int(c.B),
	}
}

func (p *XYZ) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(p.X),
		G: uint8(p.Y),
		B: uint8(p.Z),
	}
}

func (p *XYZ) DistTORGBA(c color.RGBA) float64 {
	dx := float64(c.A) - float64(p.X)
	dy := float64(c.G) - float64(p.Y)
	dz := float64(c.B) - float64(p.Z)
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

type KMeans struct {
	k     int
	cntrs []XYZ
	dist  []float64
	rnd   *rand.Rand
}

func (km *KMeans) NewKMeans(k int) *KMeans {
	return &KMeans{
		k:     k,
		cntrs: make([]XYZ, k),
		dist:  make([]float64, k),
		rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (km *KMeans) DistToCntrs(c color.RGBA) {
	for i, cntr := range km.cntrs {
		km.dist[i] = cntr.DistTORGBA(c)
	}
}

func (km *KMeans) DistArgMin() int {
	min := km.dist[0]
	min_i := 0
	for i, d := range km.dist {
		if d < min {
			min = d
			min_i = i
		}
	}
	return min_i
}

func (km *KMeans) Clusterize(colors []color.RGBA) {
	for i := range km.cntrs {
		km.cntrs[i] = RGBAToXYZ(colors[km.rnd.Intn(len(colors))])
	}

	new_cntrs := make([]XYZ, km.k)
	for {
		for i := range new_cntrs {
			new_cntrs[i].X = 0
			new_cntrs[i].Y = 0
			new_cntrs[i].Z = 0
			new_cntrs[i].C = 0
		}

		for _, color := range colors {
			km.DistToCntrs(color)
			min_i := km.DistArgMin()
			new_cntrs[min_i].X += int(color.R)
			new_cntrs[min_i].Y += int(color.G)
			new_cntrs[min_i].Z += int(color.B)
			new_cntrs[min_i].C += 1
		}

		for i := range new_cntrs {
			new_cntrs[i].X /= new_cntrs[i].C
			new_cntrs[i].Y /= new_cntrs[i].C
			new_cntrs[i].Z /= new_cntrs[i].C
		}

		is_equal := true
		for i, cntr := range new_cntrs {
			if km.cntrs[i] != cntr {
				is_equal = false
				break
			}
		}
		if is_equal {
			return
		}

		km.cntrs = new_cntrs
	}
}
