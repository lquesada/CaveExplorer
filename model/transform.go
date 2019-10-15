package model

import "github.com/lquesada/cavernal/lib/g3n/engine/core"
import "github.com/lquesada/cavernal/lib/g3n/engine/math32"

func (s *Step) ApplyTransforms(progress float32) {
	for target, t := range s.transformSteps {
		if target.Position != nil {
			if t.Mutate != nil && t.Mutate.Position != nil {
				target.Position.X = t.Set.Position.X * (1-progress) + t.Mutate.Position.X * progress
				target.Position.Y = t.Set.Position.Y * (1-progress) + t.Mutate.Position.Y * progress
				target.Position.Z = t.Set.Position.Z * (1-progress) + t.Mutate.Position.Z * progress
			} else if t.Set.Position != nil {
				target.Position.Copy(t.Set.Position)
			}
		}
		if target.Rotation != nil {
			if t.Mutate != nil && t.Mutate.Rotation != nil {
				target.Rotation.X = t.Set.Rotation.X * (1-progress) + t.Mutate.Rotation.X * progress
				target.Rotation.Y = t.Set.Rotation.Y * (1-progress) + t.Mutate.Rotation.Y * progress
				target.Rotation.Z = t.Set.Rotation.Z * (1-progress) + t.Mutate.Rotation.Z * progress
			} else if t.Set.Rotation != nil {
				target.Rotation.Copy(t.Set.Rotation)
			}
		}
		if target.Scale != nil {
			if t.Mutate != nil && t.Mutate.Scale != nil {
				target.Scale.X = t.Set.Scale.X * (1-progress) + t.Mutate.Scale.X * progress
				target.Scale.Y = t.Set.Scale.Y * (1-progress) + t.Mutate.Scale.Y * progress
				target.Scale.Z = t.Set.Scale.Z * (1-progress) + t.Mutate.Scale.Z * progress
			} else if t.Set.Scale != nil {
				target.Scale.Copy(t.Set.Scale)
			}
		}
	}
}

// --

type Transform struct{
	Position *math32.Vector3
	Rotation *math32.Vector3
	Scale *math32.Vector3
}

func NewTransform() *Transform {
	return &Transform{
		Position: &math32.Vector3{},
		Rotation: &math32.Vector3{},
		Scale: &math32.Vector3{1, 1, 1},
	}
}

func (t *Transform) Apply(node *core.Node) *core.Node {
	if t.Position != nil || t.Rotation != nil  || t.Scale != nil {
		nTmp := core.NewNode()
		nTmp.Add(node)
		if t.Position != nil && t.Position.Length() > 0 {
			nTmp.SetPositionVec(t.Position)
		}
		if t.Rotation != nil && t.Rotation.Length() > 0 {
			nTmp.SetRotationVec(t.Rotation)
		}
		if t.Scale != nil && (t.Scale.X != 1 || t.Scale.Y != 1 || t.Scale.Z != 1) {
			nTmp.SetScaleVec(t.Scale)
		}
		node = nTmp
	}
	return node
}

func (t *Transform) Inverse() *Transform {
	t2 := &Transform{
	}
	if t.Position != nil {
		t2.Position = math32.NewVector3(-t.Position.X, -t.Position.Y, -t.Position.Z)
	}
	if t.Rotation != nil {
		t2.Rotation = math32.NewVector3(-t.Rotation.X, -t.Rotation.Y, -t.Rotation.Z)
	}
	if t.Scale != nil {
		t2.Scale = math32.NewVector3(1/t.Scale.X, 1/t.Scale.Y, 1/t.Scale.Z)
	}
	return t
}

// --

type TransformStep struct{
	Set *Transform
	Mutate *Transform
}