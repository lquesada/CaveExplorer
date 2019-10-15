package humanoid

import "github.com/lquesada/cavernal/entity"


var (
	HumanoidBody entity.BodyType = "Humanoid"

	Helmet entity.ItemType = "HumanoidHelmet"
	Armor entity.ItemType = "HumanoidArmor"
	Gloves entity.ItemType = "HumanoidGloves"
	Pants entity.ItemType = "HumanoidPants"
	Boots entity.ItemType = "HumanoidBoots"
	Amulet entity.ItemType = "HumanoidAmulet"
	Ring entity.ItemType = "HumanoidRing"
	Back entity.ItemType = "HumanoidBack"

AmuletSlot entity.SlotId = "Amulet"
HelmetSlot entity.SlotId = "Helmet"
ArmorSlot entity.SlotId = "Armor"
GlovesSlot entity.SlotId = "Gloves"
PantsSlot entity.SlotId = "Pants"
BootsSlot entity.SlotId = "Boots"
RingRightSlot entity.SlotId = "RingRight"
RingLeftSlot entity.SlotId = "RingLeft"
BackSlot entity.SlotId = "Back"
HandRightSlot entity.SlotId = "HandRight"
HandLeftSlot entity.SlotId = "HandLeft"
AlternateRightSlot entity.SlotId = "AlternateRight"
AlternateLeftSlot entity.SlotId = "AlternateLeft"

Consumable1Slot entity.SlotId = "Consumable1Slot"
Consumable2Slot entity.SlotId = "Consumable2Slot"
Consumable3Slot entity.SlotId = "Consumable3Slot"
Consumable4Slot entity.SlotId = "Consumable4Slot"
Consumable5Slot entity.SlotId = "Consumable5Slot"
Consumable6Slot entity.SlotId = "Consumable6Slot"
Consumable7Slot entity.SlotId = "Consumable7Slot"
Consumable8Slot entity.SlotId = "Consumable8Slot"
Consumable9Slot entity.SlotId = "Consumable9Slot"
Consumable10Slot entity.SlotId = "Consumable10Slot"

	Head entity.DrawSlotId = "Head"
	Torso entity.DrawSlotId = "Torso"
	Face entity.DrawSlotId = "Face"
	LegLeft entity.DrawSlotId = "LegLeft"
	LegRight entity.DrawSlotId = "legRight"
	ArmLeft entity.DrawSlotId = "ArmLeft"
	ArmRight entity.DrawSlotId = "ArmRight"
	HandLeft entity.DrawSlotId = "HandLeft"
	HandRight entity.DrawSlotId = "HandRight"

ProtrudeHeadUpOrBackLong entity.DrawSlotId = "ProtrudeHeadUpOrBackLong"
ProtrudeHeadUpOrBackShort entity.DrawSlotId = "ProtrudeHeadUpOrBackShort"
ProtrudeHeadFrontOrHangLong entity.DrawSlotId = "ProtrudeHeadFrontOrHangLong"
ProtrudeHeadFrontOrHangShort entity.DrawSlotId = "ProtrudeHeadFrontOrHangShort"
ProtrudeHeadSidesLong entity.DrawSlotId = "ProtrudeHeadSidesLong"
ProtrudeHeadSidesShort entity.DrawSlotId = "ProtrudeHeadSidesShort"
ProtrudeHandLeft entity.DrawSlotId = "ProtrudeHandsLeft"
ProtrudeHandRight entity.DrawSlotId = "ProtrudeHandsRight"
ProtrudeNeck entity.DrawSlotId = "ProtrudeNeck"



StandardSlots = []entity.SlotId{
	AmuletSlot,
	HelmetSlot,
	ArmorSlot,
	GlovesSlot,
	PantsSlot,
	BootsSlot,
	RingRightSlot,
	RingLeftSlot,
	BackSlot,
	HandRightSlot,
	HandLeftSlot,
	AlternateRightSlot,
	AlternateLeftSlot,
	Consumable1Slot,
	Consumable2Slot,
	Consumable3Slot,
	Consumable4Slot,
	Consumable5Slot,
	Consumable6Slot,
	Consumable7Slot,
	Consumable8Slot,
	Consumable9Slot,
	Consumable10Slot,
}

EquippedSlots = map[entity.SlotId]bool {
	AmuletSlot: true,
	HelmetSlot: true,
	ArmorSlot: true,
	GlovesSlot: true,
	PantsSlot: true,
	BootsSlot: true,
	RingRightSlot: true,
	RingLeftSlot: true,
	BackSlot: true,
	HandRightSlot: true,
	HandLeftSlot: true,
	AlternateRightSlot: true,
	AlternateLeftSlot: true,
}

ConsumableSlots = map[entity.SlotId]bool {
	Consumable1Slot: true,
	Consumable2Slot: true,
	Consumable3Slot: true,
	Consumable4Slot: true,
	Consumable5Slot: true,
	Consumable6Slot: true,
	Consumable7Slot: true,
	Consumable8Slot: true,
	Consumable9Slot: true,
	Consumable10Slot: true,
}

ConsumableSlotsSorted = []entity.SlotId {
	Consumable1Slot,
	Consumable2Slot,
	Consumable3Slot,
	Consumable4Slot,
	Consumable5Slot,
	Consumable6Slot,
	Consumable7Slot,
	Consumable8Slot,
	Consumable9Slot,
	Consumable10Slot,
}

equipmentSlotFitItem = map[entity.SlotId]map[entity.ItemType]bool {
AmuletSlot: map[entity.ItemType]bool{Amulet: true},
HelmetSlot: map[entity.ItemType]bool{Helmet: true},
ArmorSlot: map[entity.ItemType]bool{Armor: true},
GlovesSlot: map[entity.ItemType]bool{Gloves: true},
PantsSlot: map[entity.ItemType]bool{Pants: true},
BootsSlot: map[entity.ItemType]bool{Boots: true},
RingRightSlot: map[entity.ItemType]bool{Ring: true},
RingLeftSlot: map[entity.ItemType]bool{Ring: true},
BackSlot: map[entity.ItemType]bool{Back: true},
HandRightSlot: entity.HandItemType,
HandLeftSlot: entity.HandItemType,
AlternateRightSlot: entity.HandItemType,
AlternateLeftSlot: entity.HandItemType,
}

slotsFitItem = map[entity.ItemType]entity.SlotId {
	Amulet: AmuletSlot,
	Helmet: HelmetSlot,
	Armor: ArmorSlot,
	Gloves: GlovesSlot,
	Pants: PantsSlot,
	Boots: BootsSlot,
	Back: BackSlot,
}

handToSymmetricSlot = map[entity.SlotId]entity.SlotId {
	HandRightSlot: HandLeftSlot,
	HandLeftSlot: HandRightSlot,
	AlternateRightSlot: AlternateLeftSlot,
	AlternateLeftSlot: AlternateRightSlot,
}

symmetricSlot = map[entity.SlotId]entity.SlotId {
	RingRightSlot: RingLeftSlot,
	RingLeftSlot: RingRightSlot,
	HandRightSlot: HandLeftSlot,
	HandLeftSlot: HandRightSlot,
	AlternateRightSlot: AlternateLeftSlot,
	AlternateLeftSlot: AlternateRightSlot,
}
)
