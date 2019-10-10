package consumables

import (
"cavernal.com/assets"
"cavernal.com/helpers"
"cavernal.com/model"
"cavernal.com/entity"
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
