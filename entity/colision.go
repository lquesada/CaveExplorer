package entity

import "cavernal.com/lib/g3n/engine/math32"

// RelativeCylinder is a cylinder relative to a bigger entity, e.g. a sword that's carried.
type RelativeCylinder struct{
	Ahead float32
	Y float32
	SimpleCylinder *SimpleCylinder
}

// SimpleCylinder is a bare abstract cylinder with a certain volume.
type SimpleCylinder struct{
	Radius float32
	Height float32
}

// Cylinder is a simple cylinder somewhere in the 3D space.
type Cylinder struct{
	X float32
	Y float32
	Z float32
	SimpleCylinder *SimpleCylinder
}

func GenerateCylinders(maxWidth, reach, height, y, startAhead float32) []*RelativeCylinder {
	if reach < maxWidth {
		return []*RelativeCylinder{
				&RelativeCylinder{Ahead: startAhead, Y: y, SimpleCylinder: &SimpleCylinder{Radius: reach/2, Height: height}},
			}
	}
	cylinders := []*RelativeCylinder{}
	ahead := startAhead
	for ahead < reach-maxWidth/2 {
		cylinders = append(cylinders,&RelativeCylinder{Ahead: ahead, Y: y, SimpleCylinder: &SimpleCylinder{Radius: maxWidth/2, Height: height}},)
		ahead += maxWidth/2
	}
	ahead = reach-maxWidth/2
	cylinders = append(cylinders,&RelativeCylinder{Ahead: ahead, Y: y, SimpleCylinder: &SimpleCylinder{Radius: maxWidth/2, Height: height}})
	return cylinders
}

func SimpleToAbsolute(s *SimpleCylinder) *Cylinder {
	return &Cylinder{
		X: 0,
		Y: 0,
		Z: 0,
		SimpleCylinder: s,
	}
}

func SimpleToAbsoluteList(s []*SimpleCylinder) []*Cylinder {
	a := make([]*Cylinder, len(s), len(s))
	for i, v := range s {
		a[i] = SimpleToAbsolute(v)
	}
	return a
}

func RelativeToAbsolute(r *RelativeCylinder, lookAngle float32) *Cylinder {
	return &Cylinder{
		X: r.Ahead*math32.Sin(lookAngle),
		Y: r.Y,
		Z: r.Ahead*math32.Cos(lookAngle),
		SimpleCylinder: r.SimpleCylinder,
	}
}

func RelativeToAbsoluteList(r []*RelativeCylinder, lookAngle float32) []*Cylinder {
	a := make([]*Cylinder, len(r), len(r))
	for i, v := range r {
		a[i] = RelativeToAbsolute(v, lookAngle)
	}
	return a
}

func CheckColisionBetweenFrames(x1Old, y1Old, z1Old, x1New, y1New, z1New float32, cylinder1 []*Cylinder, x2Old, y2Old, z2Old, x2New, y2New, z2New float32, cylinder2 []*Cylinder) bool {
	for _, c1 := range cylinder1 {
		for _, c2 := range cylinder2 {
			y1Min := math32.Min(y1Old, y1New)+c1.Y
			y1Max := math32.Max(y1Old, y1New)+c1.Y
			y2Min := math32.Min(y2Old, y2New)+c2.Y // NOTE c2.Y was missing, added it. may have introduced a bug
			y2Max := math32.Max(y2Old, y2New)+c2.Y // NOTE c2.Y was missing, added it. may have introduced a bug

			if math32.Max(y1Min, y2Min) <= math32.Min(y1Max + c1.SimpleCylinder.Height, y2Max + c2.SimpleCylinder.Height) {
				distance := ShortestDistanceBetweenLines(x1Old+c1.X, z1Old+c1.Z, x1New+c1.X, z1New+c1.Z, x2Old+c2.X, z2Old+c2.Z, x2New+c2.X, z2New+c2.Z)
				if distance < c1.SimpleCylinder.Radius+c2.SimpleCylinder.Radius {
					return true
				}
			}
		}
	}
	return false
}

func ShortestDistanceBetweenLines(x1Old, z1Old, x1New, z1New, x2Old, z2Old, x2New, z2New float32) float32 {
var epsilon float32 = 0.001
uX := x1New-x1Old+epsilon
uY := z1New-z1Old+epsilon

vX := x2New-x2Old+epsilon
vY := z2New-z2Old+epsilon

wX := x1Old-x2Old+epsilon
wY := z1Old-z2Old+epsilon

a := uX*uX+uY*uY
b := uX*vX+uY*vY
c := vX*vX+vY*vY
d := uX*wX+uY*wY
e := vX*wX+vY*wY

nD := a*c - b*b
var sc float32
var sN float32
sD := nD

var tc float32
var tN float32
tD := nD
        
if math32.Abs(nD) < epsilon {
  sN = 0.0
  sD = 1.0
  tN = e
  tD = c
} else {
  sN = (b*e - c*d)
  tN = (a*e - b*d)
  if sN < 0.0 {
      sN = 0.0
      tN = e
      tD = c
  } else if sN > sD {
      sN = sD
      tN = e + b
      tD = c
  }
}
        
if tN < 0.0 {
  tN = 0.0
  if -d < 0.0 {
    sN = 0.0
  } else if -d > a {
    sN = sD
  } else {
    sN = -d
    sD = a
  }
} else if tN > tD {
  tN = tD
  if (-d + b) < 0.0 {
    sN = 0
  } else if (-d + b) > a {
    sN = sD
  } else {
    sN = (-d + b)
    sD = a
  }
}

if math32.Abs(sD) < epsilon {
  sc = 0.0
} else {
  sc = sN / sD
}
if math32.Abs(tD) < epsilon {
  tc = 0.0
} else {
  tc = tN / tD
}

dpx := wX+uX*sc-vX*tc
dpz := wY+uY*sc-vY*tc

return math32.Sqrt(dpx*dpx+dpz*dpz)
}