package entity

import "cavernal.com/lib/g3n/engine/math32"

func Distance2D(v1 *math32.Vector3, v2 *math32.Vector3) float32 {
	var x = v1.X-v2.X
	var z = v1.Z-v2.Z
	return math32.Sqrt(x*x+z*z)
}
