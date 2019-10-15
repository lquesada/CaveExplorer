package customization

import "github.com/lquesada/cavernal/model"
import "github.com/lquesada/cavernal/assets"

type pathMarker struct{} // Needed to find directory
var (
	dir = model.DirOf(pathMarker{})

	FaceScaredModel = &model.NodeSpec{
		Decoder: model.Load(dir, "facescared", assets.Files),
	}
)
