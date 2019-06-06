package logger

func ExampleLogger() {
	SetName("a")
	Debugf("%d|%s", 3, "--")
	SetName("b")
	Warnf("%d|%s", 3, "--")
	SetName("")
	Infof("%d|%s", 3, "--")
	SetName("c")
	Errorf("%d|%s", 3, "--")
	//Output:
	//3|--
	//3|--
	//3|--
	//3|--
}
