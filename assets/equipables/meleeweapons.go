package equipables

import (
"cavernal.com/assets"
"cavernal.com/helpers"
"cavernal.com/entity"
"cavernal.com/model"
"cavernal.com/lib/g3n/engine/math32"
)

// --

var crystalStaffModel = &model.NodeSpec{
				Decoder: model.Load(dir, "crystalstaff", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 20, 0},
					Scale: model.X3.Scale,
				},
	}

func NewCrystalStaff() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "crystal staff",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 24,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "crystal staff melee attack",
	          MaxWidth: 0.3,
		      Reach: 1.4,
		      Height: 1.6,
    	}).New,
		EquippedNode: crystalStaffModel.Build(),
		DroppedNode: crystalStaffModel.Build(),
		InventoryNode: crystalStaffModel.Build(),
		}).New()
}

// --

var ironSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "ironsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 4.4, 0},
					Scale: model.X4.Scale,
				},
	}

func NewIronSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "iron sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 17,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "iron sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: ironSwordModel.Build(),
		DroppedNode: ironSwordModel.Build(),
		InventoryNode: ironSwordModel.Build(),
		}).New()
}

// --

var grassSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "grasssword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 4, 0},
					Scale: model.X4.Scale,
				},
	}

func NewGrassSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "grass sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 33,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "grass sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: grassSwordModel.Build(),
		DroppedNode: grassSwordModel.Build(),
		InventoryNode: grassSwordModel.Build(),
		}).New()
}

// --

var finnSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "finnsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 3, 0},
					Scale: model.X4.Scale,
				},
	}

func NewFinnSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "Finn's sword",
		ItemType: entity.TwoHandedWeapon,
		AttackValue: 37,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "Finn's sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: finnSwordModel.Build(),
		DroppedNode: finnSwordModel.Build(),
		InventoryNode: finnSwordModel.Build(),
		}).New()
}

// --

var demonbloodSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "demonbloodsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 6, 0},
					Scale: model.X4.Scale,
				},
	}

func NewDemonbloodSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "demonblood sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 47,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "demonblood sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: demonbloodSwordModel.Build(),
		DroppedNode: demonbloodSwordModel.Build(),
		InventoryNode: demonbloodSwordModel.Build(),
		}).New()
}

// --

var longSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "longsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 9, 0},
					Scale: model.X2.Scale,
				},
	}

func NewLongSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "long sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 22,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "long sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: longSwordModel.Build(),
		DroppedNode: longSwordModel.Build(),
		InventoryNode: longSwordModel.Build(),
		}).New()
}

// --

var rapierModel = &model.NodeSpec{
				Decoder: model.Load(dir, "rapier", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 4, 0},
					Scale: model.X4.Scale,
				},
	}

func NewRapier() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Common,
		Name: "rapier",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 16,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "rapier melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: rapierModel.Build(),
		DroppedNode: rapierModel.Build(),
		InventoryNode: rapierModel.Build(),
		}).New()
}


// --

var rootSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "rootsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 3, 0},
					Scale: model.X4.Scale,
				},
	}

func NewRootSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "root sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 28,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "root sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: rootSwordModel.Build(),
		DroppedNode: rootSwordModel.Build(),
		InventoryNode: rootSwordModel.Build(),
		}).New()
}

// --

var scarletSwordModel = &model.NodeSpec{
				Decoder: model.Load(dir, "scarletsword", assets.Files),
				Transform: &model.Transform{
					Position: &math32.Vector3{0, 5, 0},
					Scale: model.X4.Scale,
				},
	}

func NewScarletSword() *helpers.BasicEquipable {
	return (&helpers.BasicEquipableSpecification{
		Category: entity.Artifact,
		Name: "scarlet sword",
		ItemType: entity.OneHandedWeapon,
		AttackValue: 32,
		AttackGenerator: (&helpers.BasicMeleeAttackSpecification{
		      Name: "scarlet sword melee attack",
	          MaxWidth: 0.3,
		      Reach: 1,
		      Height: 1.6,
    	}).New,
		EquippedNode: scarletSwordModel.Build(),
		DroppedNode: scarletSwordModel.Build(),
		InventoryNode: scarletSwordModel.Build(),
		}).New()
}

// --
