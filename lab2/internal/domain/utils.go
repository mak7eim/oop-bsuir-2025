package domain

func SearchByID[T Entity](entities []T, id string) (T, bool) {
	for _, e := range entities {
		if e.GetID() == id{
			return e, true
		}
	}
	var zero T
	return zero, false
}