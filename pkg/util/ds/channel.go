package ds

func ReadAll[T any](ch <-chan T) []T {
	var items []T
read:
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				break read
			}
			items = append(items, item)
		default:
			break read
		}
	}
	return items
}
