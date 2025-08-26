package eval

import "errors"

var errNotSupportedOnWindows = errors.New("not supported on Windows")

func execFn(...any) error {
	return errNotSupportedOnWindows
}

// fg implements foreground job control on Windows using Job Objects.
func fg(pids ...int) error {
	controller, err := NewJobController()
	if err != nil {
		return err
	}
	defer controller.Close()

	return fgWindows(controller, pids...)
}
