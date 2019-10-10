package customization

import "cavernal.com/model"
import "cavernal.com/assets"

type pathMarker struct{} // Needed to find directory
var (
	dir = model.DirOf(pathMarker{})

	FaceScaredModel = &model.NodeSpec{
		Decoder: model.Load(dir, "facescared", assets.Files),
	}
)
