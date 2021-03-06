package platforms

import "github.com/dena/devfarm/cmd/core/testutil"

func AnyDeviceFinder() DeviceFinder {
	return StubDeviceFinder(false, testutil.AnyError)
}

func StubDeviceFinder(b bool, err error) DeviceFinder {
	return func(EitherDevice) (bool, error) {
		return b, err
	}
}
