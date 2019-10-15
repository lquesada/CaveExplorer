package entity

import (
	"math"
	"github.com/lquesada/cavernal/lib/g3n/engine/math32"
)

func NormalizeAngle(a float32) float32 {
	rad := float32(math.Remainder(float64(a), 2*math32.Pi))
	if rad <= -math32.Pi {
		rad = math32.Pi
	}
	return rad
}
