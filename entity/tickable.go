package entity

type ITickable interface {
	PreTick()
	Tick(delta float32)
	PostTick(delta float32)
}