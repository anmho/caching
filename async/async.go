package async

func HandleAsync(callback func()) {
	// should wrap in a panic handler
	go func() {
		callback()
	}()
}
