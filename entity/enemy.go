package entity

type IEnemy interface {
	ICharacter
	IEvil

	TargetSeeDistance() float32

	Target() IEntity

	SeeTarget(IEntity)
}

type IEvil interface {
	Evil()
}

func FilterEnemy(v interface{}) bool {
	_, ok := v.(IEnemy)
	return ok
}


func FilterEvil(v interface{}) bool {
	_, ok := v.(IEvil)
	return ok
}