package mutiny

type AssignedValue string

const (
	Fail      AssignedValue = "FAIL"
	Erroneous AssignedValue = "ERRONEOUS"
	Nil       AssignedValue = "NIL"
)

type PossibleValues struct {
	use       AssignedValue
	Pass      []any
	Fail      []any
	Erroneous []any
}

func SelectValue(pv PossibleValues, use AssignedValue) PossibleValues {
	pv.use = use
	return pv
}
