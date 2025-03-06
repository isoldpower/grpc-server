package util

func ProtectedAction(
	err error,
	action func() error,
) error {
	if err != nil {
		return err
	}

	return action()
}
