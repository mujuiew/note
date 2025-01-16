package null

type ExistableField interface {
	Null() bool
	EqualsI(interface{}) (bool, error)
}
type ComparableField interface {
	ExistableField
	GTI(interface{}) (bool, error)
	GTEI(interface{}) (bool, error)
	LTI(interface{}) (bool, error)
	LTEI(interface{}) (bool, error)
}
type LikeableField interface {
	LikeI(interface{}) (bool, error)
}
