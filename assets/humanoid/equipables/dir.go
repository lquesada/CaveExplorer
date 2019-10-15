package equipables

import "github.com/lquesada/cavernal/model"

type pathMarker struct{} // Needed to find directory
var dir = model.DirOf(pathMarker{})
