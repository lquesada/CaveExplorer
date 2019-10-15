package consumables

import (
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/helpers"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/entity"
)

type pathMarker struct{} // Needed to find directory
var (
	dir = model.DirOf(pathMarker{})
)

// --

var bananaModel = &model.NodeSpec{
				Decoder: model.Load(dir, "banana", assets.Files),
				Transform: model.X2,
	}

var bananaCountable entity.CountableId = "banana"

func NewBanana() *helpers.BasicConsumable {
	return (&helpers.BasicConsumableSpecification{
		Category: entity.Common,
		Name: "banana",
		DroppedNode: bananaModel.Build(),
		InventoryNode: bananaModel.Build(),
		CountableId: bananaCountable,
		}).New()
}

// --
