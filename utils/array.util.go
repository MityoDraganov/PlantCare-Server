package utils

// ReturnEmptyIfNil checks if a slice is nil or empty, and returns an empty slice if so
func ReturnEmptyIfNil[T any](data []T) []T {
    if data == nil {
        return []T{} // Return an empty array if it's nil
    }
    return data
}
