package async

func HandleAsync(callback func()) {
	go func() {
		callback()
	}()
}
