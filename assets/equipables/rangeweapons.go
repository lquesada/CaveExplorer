package equipables

import (
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/helpers"
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/lib/g3n/engine/math32"
)

// --

var shortBowModel = &model.NodeSpec{
				Decoder: model.Load(dir, "shortbow", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 12, 6},
					Scale: model.X3.Scale,
				},
	}

var arrowAmmo entity.AmmoType = "arrow"

func NewShortBow() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "short bow",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 5,
		EquippedNode: shortBowModel.Build(),
		DroppedNode: shortBowModel.Build(),
		InventoryNode: shortBowModel.Build(),
		RequiresAmmo: arrowAmmo,
		}).New()
}

// --

var arrowModel = &model.NodeSpec{
				Decoder: model.Load(dir, "arrow", assets.Files),
				Transform: &model.Transform{
					Scale: model.X2.Scale,
				},
	}

var arrowCountable entity.CountableId = "arrow"


func NewArrow() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "arrow",
		ItemType: entity.Other,
		AttackValue: 5,
		EquippedNode: arrowModel.Build(),
		DroppedNode: arrowModel.Build(),
		InventoryNode: arrowModel.Build(),
		CountableId: arrowCountable,
		ProvidesAmmo: arrowAmmo,
		}).New()
}
