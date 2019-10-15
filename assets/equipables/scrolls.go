package equipables

import (
"github.com/lquesada/cavernal/assets"
"github.com/lquesada/cavernal/helpers"
"github.com/lquesada/cavernal/entity"
"github.com/lquesada/cavernal/model"
"github.com/lquesada/cavernal/lib/g3n/engine/math32"
)

// --

var fireScrollModel = &model.NodeSpec{
				Decoder: model.Load(dir, "firescroll", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 12, 6},
					Scale: model.X2.Scale,
				},
	}

var fireScrollCountable entity.CountableId = "fireScroll"

var fireScrollAmmo entity.AmmoType = "fireScroll"

func NewFireScroll() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Rare,
		Name: "fire scroll",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 20,
		EquippedNode: fireScrollModel.Build().Transform(&model.Transform{
								Rotation: &math32.Vector3{0, math32.Pi, 0},
							}),
		DroppedNode: fireScrollModel.Build(),
		InventoryNode: fireScrollModel.Build(),
		CountableId: fireScrollCountable,
		RequiresAmmo: fireScrollAmmo,
		ProvidesAmmo: fireScrollAmmo,
		}).New()
}

// --
