package util

import "sync"

func ProtectedAction(
	err error,
	action func() error,
) error {
	if err != nil {
		return err
	}

	return action()
}

func transferRoutine[T any](done <-chan bool, waitGroup *sync.WaitGroup, inputStream <-chan T, outputStream chan<- T) {
	defer waitGroup.Done()
	for value := range inputStream {
		select {
		case <-done:
			return
		case outputStream <- value:
			break
		}
	}
}

func FlatStreams[T any](done <-chan bool, streams ...<-chan T) <-chan T {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(streams))
	outputStream := make(chan T)

	for _, stream := range streams {
		go transferRoutine(done, waitGroup, stream, outputStream)
	}

	go func() {
		waitGroup.Wait()
		close(outputStream)
	}()

	return outputStream
}
