package model

import (
		"bytes"
		"reflect"
		"strings"
        "fmt"

        "github.com/lquesada/cavernal/lib/g3n/engine/loader/obj"
        "github.com/lquesada/cavernal/lib/g3n/engine/math32"
)

var (
MeterScaleFactor float32 = 1/1.26
GlobalScaleFactor float32 = 0.01
X2 = &Transform{Scale: &math32.Vector3{2, 2, 2}}
X3 = &Transform{Scale: &math32.Vector3{3, 3, 3}}
X4 = &Transform{Scale: &math32.Vector3{4, 4, 4}}
X5 = &Transform{Scale: &math32.Vector3{5, 5, 5}}
X6 = &Transform{Scale: &math32.Vector3{6, 6, 6}}
X7 = &Transform{Scale: &math32.Vector3{7, 7, 7}}
X8 = &Transform{Scale: &math32.Vector3{8, 8, 8}}
Xhalf = &Transform{Scale: &math32.Vector3{0.5, 0.5, 0.5}}
Xthird = &Transform{Scale: &math32.Vector3{0.333, 0.333, 0.333}}
Xquarter = &Transform{Scale: &math32.Vector3{0.25, 0.25, 0.25}}
)

func DirOf(i interface{}) string {
	return strings.TrimPrefix(reflect.TypeOf(i).PkgPath(), "github.com/lquesada/cavernal/")
}

func load(filepath, baseFilename string, files map[string][]byte) (*obj.Decoder, error) {
	objFilename := fmt.Sprintf("%s/%s.obj", filepath, baseFilename)
	mtlFilename := fmt.Sprintf("%s/%s.mtl", filepath, baseFilename)
	objData, ok := files[objFilename]
	if !ok {
        return nil, fmt.Errorf("obj file %s/%s data not found", filepath, baseFilename)
	}
	mtlData, ok := files[mtlFilename]
	if !ok {
        return nil, fmt.Errorf("mtl file %s/%s data not found", filepath, baseFilename)
	}
        dec, err := obj.DecodeReader(bytes.NewReader(objData), bytes.NewReader(mtlData), filepath, files)
        if err != nil {
                return nil, fmt.Errorf("couldn't decode object: %s", err)
        }
        _, err = dec.NewGroup()
        if err != nil {
                return nil, fmt.Errorf("couldn't get group for %s/%s{.obj/.mtl/.png}: %s", filepath, baseFilename, err)
        }
        return dec, nil
}

func Load(filepath, baseFilename string, files map[string][]byte) *obj.Decoder {
        model, err := load(filepath, baseFilename, files)
        if err != nil {
                panic(fmt.Sprintf("Critical error. Failed to load %s/%s{.obj/.mtl/.png}: %s", filepath, baseFilename, err))
        }
	    _, err = model.NewGroup()
	    if err != nil {
            	panic(fmt.Sprintf("Critical error. Failed to get test group for %s/%s{.obj/.mtl/.png}: %s", filepath, baseFilename, err))
    	}
        return model
}

type NodeSpec struct {
	Decoder *obj.Decoder
	Transform *Transform
	RGBA *math32.Color4
}

func (m *NodeSpec) Build() *node {
    n, err := m.Decoder.NewGroup()
    if err != nil {
            panic(fmt.Sprintf("Critical error. Failed to get group: %s", err))
    }
	position := &math32.Vector3{}
	if m.Transform != nil && m.Transform.Position != nil {
		position.Copy(m.Transform.Position)
		position.MultiplyScalar(-GlobalScaleFactor)
		if m.Transform != nil && m.Transform.Scale != nil && m.Transform.Scale.Length() != 0 {
			position.Multiply(m.Transform.Scale)
		}
	}
	rotation := &math32.Vector3{}
	if m.Transform != nil && m.Transform.Rotation != nil {
	rotation.Copy(m.Transform.Rotation)
	}
	scale := &math32.Vector3{GlobalScaleFactor, GlobalScaleFactor, GlobalScaleFactor}
	if m.Transform != nil && m.Transform.Scale != nil && m.Transform.Scale.Length() != 0 {
		scale.Multiply(m.Transform.Scale)
	}
	rgba := &math32.Color4{1, 1, 1, 1}
	if m.RGBA != nil {
		rgba.R, rgba.G, rgba.B, rgba.A = m.RGBA.R, m.RGBA.G, m.RGBA.B, m.RGBA.A
	}
	return &node{
		nodes: []INode{&modelNode{node: n}},
		transform: &Transform{
			Position: position,
			Rotation: rotation,
			Scale: scale,
		},
		rgba: rgba,
	}
}