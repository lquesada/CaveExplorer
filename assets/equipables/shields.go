package equipables

import (
"cavernal.com/assets"
"cavernal.com/helpers"
"cavernal.com/entity"
"cavernal.com/model"
"cavernal.com/lib/g3n/engine/math32"
)

// --

var woodenShieldModel = &model.NodeSpec{
				Decoder: model.Load(dir, "woodenshield", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{1, 11, -1},
					Scale: model.X3.Scale,
				},
	}

func NewWoodenShield() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "wooden shield",
		ItemType: entity.Shield,
		DefenseValue: 4,
		EquippedNode: woodenShieldModel.Build(),
		DroppedNode: woodenShieldModel.Build(),
		InventoryNode: woodenShieldModel.Build(),
		}).New()
}

// --

var reinforcedShieldModel = &model.NodeSpec{
				Decoder: model.Load(dir, "reinforcedshield", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{1, 11, -1},
					Scale: model.X3.Scale,
				},
	}

func NewReinforcedShield() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "reinforced shield",
		ItemType: entity.Shield,
		DefenseValue: 5,
		EquippedNode: reinforcedShieldModel.Build(),
		DroppedNode: reinforcedShieldModel.Build(),
		InventoryNode: reinforcedShieldModel.Build(),
		}).New()
}

// --

var ironShieldModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironshield", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{1, 11, -1},
					Scale: model.X3.Scale,
				},
	}

func NewIronShield() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "iron shield",
		ItemType: entity.Shield,
		DefenseValue: 7,
		EquippedNode: ironShieldModel.Build(),
		DroppedNode: ironShieldModel.Build(),
		InventoryNode: ironShieldModel.Build(),
		}).New()
}

// --

var spikeShieldModel = &model.NodeSpec{
				Decoder: model.Load(dir, "spikeshield", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{1, 11, -1},
					Scale: model.X3.Scale,
				},
	}

func NewSpikeShield() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "spike shield",
		ItemType: entity.Shield,
		DefenseValue: 8,
		EquippedNode: spikeShieldModel.Build(),
		DroppedNode: spikeShieldModel.Build(),
		InventoryNode: spikeShieldModel.Build(),
		}).New()
}

// --
