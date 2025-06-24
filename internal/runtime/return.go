package runtime

type Return struct {
	Value any
}

func (r Return) Error() string {
	if r.Value == nil {
		return "Return called with no value"
	}
	return "Return called with value: " + r.Value.(string)
}
