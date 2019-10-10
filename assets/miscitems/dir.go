package miscitems

import "cavernal.com/model"

type pathMarker struct{} // Needed to find directory
var dir = model.DirOf(pathMarker{})
