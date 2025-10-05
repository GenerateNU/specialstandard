package testutil

// Ptr returns a pointer to the given value.
// For creating pointers to literals in tests.
//
// Example:
//
//	grade := testutil.Ptr(5)
//	name := testutil.Ptr("John Doe")
//
// instead of using helper functions like:
//
//	func ptrInt(v int) *int { return &v }
func Ptr[T any](v T) *T {
	return &v
}
