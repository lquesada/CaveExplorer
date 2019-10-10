package entity

type IPlayer interface {
	ICharacter
	IGood
}

type IGood interface {
	Good()
}

func FilterPlayer(v interface{}) bool {
	_, ok := v.(IPlayer)
	return ok
}

func FilterGood(v interface{}) bool {
	_, ok := v.(IGood)
	return ok
}