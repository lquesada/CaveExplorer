package entity

type Filter func(interface{}) bool

type BodyType string

type AmmoType string

type CountableId string

type DrawSlotId string

type ItemType string

type SlotId string

type ItemCategory struct{}

var Common = ItemCategory{}
var Rare = ItemCategory{}
var Artifact = ItemCategory{}