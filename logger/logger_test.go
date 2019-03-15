package logger

func ExampleLogger() {
	//SetName("")
	Debugf("%d|%s", 3, "--")
	Warnf("%d|%s", 3, "--")
	Infof("%d|%s", 3, "--")
	Errorf("%d|%s", 3, "--")
	//Output: 3
}
