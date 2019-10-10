package miscitems

import (
"cavernal.com/assets"
"cavernal.com/helpers"
"cavernal.com/model"
"cavernal.com/entity"
)

// --

var goldNuggetModel = &model.NodeSpec{
		Decoder: model.Load(dir, "goldnugget", assets.Files),
		Transform: model.Xhalf,
	}

var goldNuggetCountable entity.CountableId = "goldNugget"

func NewGoldNugget() *helpers.BasicItem {
	return (&helpers.BasicItemSpecification{
		Category: entity.Rare,
		Name: "gold nugget",
		DroppedNode: goldNuggetModel.Build(),
		InventoryNode: goldNuggetModel.Build(),
		CountableId: goldNuggetCountable,
		}).New()
}

// --
