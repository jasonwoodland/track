package cleanup

var Cleanup func()

func SetCleanupFn(fn func()) {
	Cleanup = fn
}
