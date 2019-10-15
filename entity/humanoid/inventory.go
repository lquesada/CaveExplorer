package humanoid

import (
	"github.com/lquesada/cavernal/entity"
)

type Inventory struct {
	slots map[entity.SlotId]entity.IItem
	items []entity.IItem
	activeSetup int
	lastRingAutoEquipped int

	holdItem entity.IItem
	holdAmount int

	handsBusy bool

	generate []entity.IItem
}

func NewInventory(equippedSlots []entity.SlotId, itemCount int) *Inventory {
	slots := map[entity.SlotId]entity.IItem{}
	for _, v := range equippedSlots {
		slots[v] = nil
	}
	return &Inventory{
		slots: slots,
		items: make([]entity.IItem, itemCount),
		generate: []entity.IItem{},
	}
}

func (i *Inventory) PreTick() {
	for _, e := range i.AllItems() {
		if e != nil {
			e.PreTick()
		}
	}
}

func (i *Inventory) Tick(delta float32) {
	for _, e := range i.AllItems() {
		if e != nil {
			e.Tick(delta)
		}
	}
}

func (i *Inventory) PostTick(delta float32) {
	for _, e := range i.AllItems() {
		if e != nil {
			e.PostTick(delta)
		}
	}
}

func (i *Inventory) EquipmentSlots() map[entity.SlotId]entity.IEquipable {
	res := map[entity.SlotId]entity.IEquipable{}
	for k, v := range i.slots {
		if EquippedSlots[k] {		
			if v == nil {
				res[k] = nil
			} else {
				res[k] = v.(entity.IEquipable) 
			}
		}
	}
	return res
}

func (i *Inventory) Slots() map[entity.SlotId]entity.IItem {
	return i.slots
}

func (i *Inventory) CanEquip(itemSlot int, slot entity.SlotId) bool {
	return i.CanEquipItem(i.items[itemSlot], slot)
}

func (i *Inventory) CanEquipItem(item entity.IItem, slot entity.SlotId) bool {
	if item == nil {
		return false
	}
	if (slot == HandLeftSlot || slot == HandRightSlot) && i.handsBusy {
		return false
	}

	if i.slots[slot] != nil && item.IsCountable() && i.slots[slot].IsCountable() && item.CountableId() == i.slots[slot].CountableId() {
		return true
	}

	if _, ok := item.(entity.IConsumable); ok {
		if ConsumableSlots[slot] {
			return true
		} else {
			return false
		}
	}

	if v, ok := item.(entity.IEquipable); ok {
		if !equipmentSlotFitItem[slot][v.ItemType()] {
			return false
		}
		if entity.HandItemType[v.ItemType()] {
			otherSlot, isHandSlot := handToSymmetricSlot[slot]
			if !isHandSlot {
				return false
			}
			other := i.slots[otherSlot]
			if other == nil {
				return true
			}
			if v == other {
				return true
			}
			if v.ItemType() == entity.TwoHandedWeapon || other.ItemType() == entity.TwoHandedWeapon {
				return false
			}
			if v.ItemType() == entity.OneHandedWeapon && other.ItemType() == entity.OneHandedWeapon {
				return false
			}
			if v.ItemType() == entity.Shield && other.ItemType() == entity.Shield {
				return false
			}
			if v.ItemType() == entity.OneHandedWeapon && v.RequiresAmmo() != "" && other.RequiresAmmo() != "" && other.ProvidesAmmo() != v.RequiresAmmo() {
				return false
			}
			if other.ItemType() == entity.OneHandedWeapon && other.RequiresAmmo() != "" && v.RequiresAmmo() != "" && v.ProvidesAmmo() != other.RequiresAmmo() {
				return false
			}
			return true
		} else {
			return true
		}
	}
	return false
}
func (i *Inventory) AutoEquip(itemSlot int) {
	item := i.items[itemSlot]
	if item == nil {
		return
	}
	if item.IsCountable() {
		for slot, v := range i.slots {
			if v != nil && v.CountableId() == item.CountableId() {
				entity.Stack(item, i.slots[slot])
				i.items[itemSlot] = nil
				return
			}
		}
	}

	if _, ok := item.(entity.IConsumable); ok {
		for _, k := range ConsumableSlotsSorted {
			if i.CanEquip(itemSlot, k) {
				i.Equip(itemSlot, k)
				return
			}
		}
	}

	if slot, ok := slotsFitItem[item.ItemType()]; ok {
		if i.CanEquip(itemSlot, slot) {
			i.Equip(itemSlot, slot)
			return
		}
	}

	if item.ItemType() == Ring {
		var targetSlot entity.SlotId
		if i.slots[RingLeftSlot] == nil {
			targetSlot = RingLeftSlot
		} else if i.slots[RingRightSlot] == nil {
			targetSlot = RingRightSlot
		} else if i.lastRingAutoEquipped == 0 {
			targetSlot = RingRightSlot
			i.lastRingAutoEquipped = 1
		} else {
			targetSlot = RingLeftSlot
			i.lastRingAutoEquipped = 0
		}
		i.slots[targetSlot], i.items[itemSlot] = i.items[itemSlot], i.slots[targetSlot]
		return
	}

	if entity.HandItemType[item.ItemType()] {
		if i.handsBusy {
			return
		}
		var first, second entity.SlotId
		if entity.WeaponItemType[item.ItemType()] {
			first = HandRightSlot
			second = HandLeftSlot
		} else {
			first = HandLeftSlot
			second = HandRightSlot
		}
		if i.slots[first] == nil && i.CanEquip(itemSlot, first) {
			i.Equip(itemSlot, first)
		} else if i.slots[second] == nil && i.CanEquip(itemSlot, second) {
			i.Equip(itemSlot, second)
		} else if i.slots[first] != nil && entity.WeaponItemType[i.slots[first].ItemType()] == entity.WeaponItemType[item.ItemType()] {
			i.Equip(itemSlot, first)
		} else if  i.slots[second] != nil && entity.WeaponItemType[i.slots[second].ItemType()] == entity.WeaponItemType[item.ItemType()] {
			i.Equip(itemSlot, second)
		} else if i.CanEquip(itemSlot, first) {
			i.Equip(itemSlot, first)
		} else if i.CanEquip(itemSlot, second) {
			i.Equip(itemSlot, second)
		} else {
			slot, otherSlot := first, second
			aux := i.slots[otherSlot]
			i.slots[otherSlot] = nil
			if i.CanPickUp(aux) {
				i.PickUp(aux)
				i.Equip(itemSlot, slot)
			} else {
				i.slots[otherSlot] = aux
			}
		}
	}
}

func (i *Inventory) Equip(itemSlot int, slot entity.SlotId) {
	if i.items[itemSlot] == nil {
		return
	}
	if i.slots[slot] != nil && entity.CanStack(i.items[itemSlot], i.slots[slot]) {
		entity.Stack(i.items[itemSlot], i.slots[slot])
		i.items[itemSlot] = nil
		return
	}
	if i.CanEquip(itemSlot, slot) {
		i.slots[slot], i.items[itemSlot] = i.items[itemSlot], i.slots[slot]
	}
}

func (i *Inventory) AutoUnequip(slot entity.SlotId) {
	item := i.slots[slot]
	if (slot == HandRightSlot || slot == HandLeftSlot) && i.handsBusy {
		return
	}

	if item == nil {
		return
	}
	if where := i.whereCanPickUp(item); where != -1 {
		i.items[where] = item
		i.slots[slot] = nil
	}
}

func (i *Inventory) Unequip(slot entity.SlotId, itemSlot int) {
	item := i.slots[slot]
	if (slot == HandLeftSlot || slot == HandRightSlot) && i.handsBusy {
		return
	}
	if item == nil {
		return
	}
	if i.items[itemSlot] != nil && entity.CanStack(i.slots[slot], i.items[itemSlot]) {
		entity.Stack(i.slots[slot], i.items[itemSlot])
		i.slots[slot] = nil
		return
	}
	if i.items[itemSlot] == nil || i.CanEquip(itemSlot, slot) {
		i.items[itemSlot], i.slots[slot] = i.slots[slot], i.items[itemSlot]
	} else {
		i.AutoUnequip(slot)
	}
}

func (i *Inventory) ShuffleEquipment(slot1, slot2 entity.SlotId) bool {
	if (slot1 == HandLeftSlot || slot1 == HandRightSlot || slot2 == HandLeftSlot || slot2 == HandRightSlot) && i.handsBusy {
		return false
	}
	if slot1 == slot2 {
		return false
	}
	if i.slots[slot1] != nil && i.slots[slot2] != nil && entity.CanStack(i.slots[slot1], i.slots[slot2]) {
		entity.Stack(i.slots[slot1], i.slots[slot2])
		i.slots[slot1] = nil
		return true
	}
	if i.slots[slot1] == nil ||
	   i.slots[slot2] == nil ||
	   slot1 == symmetricSlot[slot2] ||
	   (ConsumableSlots[slot1] && ConsumableSlots[slot2]) ||
	   (i.CanEquipItem(i.slots[slot1], slot2) && i.CanEquipItem(i.slots[slot2], slot1)) {
		i.slots[slot1], i.slots[slot2] = i.slots[slot2], i.slots[slot1]
		return true
	}
	return false
}

func (i *Inventory) Items() []entity.IItem {
	return i.items
}

func (i *Inventory) AllItems() []entity.IItem {
	l := len(i.items) + len(i.slots)
	allItems := make([]entity.IItem, 0, l)
	for _, v := range i.items {
		if v != nil {
			allItems = append(allItems, v)
		}
	}
	for _, v := range i.slots {
		if v != nil {
			allItems = append(allItems, v)
		}
	}
	return allItems
}

func (i *Inventory) EquippedWeapon() entity.IEquipable {
	if v, ok := i.slots[HandRightSlot]; ok && v != nil && entity.WeaponItemType[v.ItemType()] {
		return v.(entity.IEquipable)
	} else if v, ok := i.slots[HandLeftSlot]; ok && v != nil && entity.WeaponItemType[v.ItemType()] {
		return v.(entity.IEquipable)
	}
	// TODO add default weapon
	return nil
}

func (i *Inventory) WhereHas(item entity.IItem) int {
	for k, e := range i.Items() {
		if e == item {
			return k
		}
	}
	return -1
}

func (i *Inventory) WhereEquipped(item entity.IItem) entity.SlotId {
	for k, e := range i.Slots() {
		if e == item {
			return k
		}
	}
	return ""
}
func (i *Inventory) whereCanPickUp(item entity.IItem) int {
	if item.IsCountable() {
		for k, v := range i.items {
			if v != nil && v.CountableId() == item.CountableId() {
				return k
			}
		}
	}
	for k, v := range i.items {
		if v == nil {
			return k
		}
	}
	return -1
}

func (i *Inventory) CanPickUp(item entity.IItem) bool {
	if i.whereCanPickUp(item) != -1 {
		return true
	}
	return false
}

func (i *Inventory) PickUp(item entity.IItem) int {
	if item.IsCountable() {
		for slot, v := range i.slots {
			if v != nil && v.CountableId() == item.CountableId() {
				entity.Stack(item, i.slots[slot])
				return -1
			}
		}
	}
	itemSlot := i.whereCanPickUp(item)
	if itemSlot == -1 {
		return -1
	}
	if i.items[itemSlot] != nil && i.items[itemSlot].CountableId() == item.CountableId() {
		entity.Stack(item, i.items[itemSlot])
	} else {
		i.items[itemSlot] = item
	}
	return itemSlot
}

func (i *Inventory) ShuffleItem(itemSlot1, itemSlot2 int) {
	if itemSlot1 == itemSlot2 {
		return
	}
	if i.items[itemSlot1] != nil && i.items[itemSlot2] != nil && entity.CanStack(i.items[itemSlot1], i.items[itemSlot2]) {
		entity.Stack(i.items[itemSlot1], i.items[itemSlot2])
		i.items[itemSlot1] = nil
		return
	}
	i.items[itemSlot1], i.items[itemSlot2] = i.items[itemSlot2], i.items[itemSlot1]
}

func (i *Inventory) SetHoldAmount(v int) {
	i.holdAmount = v
}

func (i *Inventory) HoldAmount() int {
	return i.holdAmount
}

func (i *Inventory) SetHoldItem(item entity.IItem) {
	i.holdItem = item 
}

func (i *Inventory) HoldItem() entity.IItem {
	if i.WhereHas(i.holdItem) == -1 && i.WhereEquipped(i.holdItem) == "" {
		i.holdItem = nil
	}
	if i.holdItem != nil && i.holdItem.IsCountable() && i.holdItem.Count() < i.holdAmount {
		i.holdItem = nil
	}
	return i.holdItem
}

func (i *Inventory) RemoveItemSlot(itemSlot int) {
		i.items[itemSlot] = nil
}
	
func (i *Inventory) RemoveEquipped(slot entity.SlotId) {
	if (slot == HandLeftSlot || slot == HandRightSlot) && i.handsBusy {
		return
	}
	i.slots[slot] = nil
}
	
func (i *Inventory) DropItemSlot(itemSlot int) {
	item := i.items[itemSlot]
	if item != nil {
		i.generate = append(i.generate, item)
		i.RemoveItemSlot(itemSlot)
	}
}

func (i *Inventory) DropEquipped(slot entity.SlotId) {
	if (slot == HandLeftSlot || slot == HandRightSlot) && i.handsBusy {
		return
	}
	item := i.slots[slot]
	if item != nil {
		i.generate = append(i.generate, item)
		i.RemoveEquipped(slot)
	}
}

func (i *Inventory) DropAllLoot() {
	for k, _ := range i.items {
		i.DropItemSlot(k)
	}
	for k := range i.slots {
		// Force drop everything by making hands not busy
		i.handsBusy = false
		i.DropEquipped(k)
	}
}

func (i *Inventory) Generate() []entity.IItem {
	list := []entity.IItem{}
	aux := i.generate
	i.generate = list
	return aux
}

func (i *Inventory) Flip() {
	if i.handsBusy {
		return
	}
	if i.slots[AlternateRightSlot] == nil && i.CanEquipItem(i.slots[AlternateLeftSlot], HandLeftSlot) {
		i.slots[HandLeftSlot], i.slots[AlternateLeftSlot] = i.slots[AlternateLeftSlot], i.slots[HandLeftSlot]
	} else	if i.slots[AlternateLeftSlot] == nil && i.CanEquipItem(i.slots[AlternateRightSlot], HandRightSlot) {
		i.slots[HandRightSlot], i.slots[AlternateRightSlot] = i.slots[AlternateRightSlot], i.slots[HandRightSlot]
	} else {
		i.slots[HandRightSlot], i.slots[AlternateRightSlot] = i.slots[AlternateRightSlot], i.slots[HandRightSlot]
		i.slots[HandLeftSlot], i.slots[AlternateLeftSlot] = i.slots[AlternateLeftSlot], i.slots[HandLeftSlot] 
	}
}

func (i *Inventory) Attack(repeat bool) interface{} {
	if i.handsBusy {
		return nil
	}
	if i.EquippedWeapon() != nil && !i.EquippedWeapon().AllowRepeat() && repeat {
		return nil
	}
	rightHand := false
	leftItem, ok := i.slots[HandLeftSlot].(entity.IEquipable)
	if !ok {
		leftItem = nil
	}
	rightItem, ok := i.slots[HandRightSlot].(entity.IEquipable)
	if !ok {
		rightItem = nil
	}
	var mainItem, otherItem entity.IEquipable
	if leftItem != nil && entity.WeaponItemType[leftItem.ItemType()] {
		mainItem, otherItem = leftItem, rightItem
	} else if rightItem != nil && entity.WeaponItemType[rightItem.ItemType()] {
		mainItem, otherItem = rightItem, leftItem
		rightHand = true
	}
	if mainItem == nil {
		return nil
	}
	if mainItem.RequiresAmmo() != "" {
		if otherItem != nil && otherItem.ProvidesAmmo() == mainItem.RequiresAmmo() {
			if otherItem.IsCountable() {
				if otherItem.Count() < 1 {
					return nil
				}
				otherItem.SetCount(otherItem.Count()-1)
			}
		} else if mainItem.ProvidesAmmo() == mainItem.RequiresAmmo() {
			if mainItem.IsCountable() {
				if mainItem.Count() < 1 {
					return nil
				}
				mainItem.SetCount(mainItem.Count()-1)
			}
		} else {
			return nil
		}
	}
	i.handsBusy = true
	return &attackInfo{
		rightHand: rightHand,
		attackFinishedCallback: func() { i.handsBusy = false },
	}
}



type attackInfo struct {
	rightHand bool
	attackFinishedCallback func()
}
