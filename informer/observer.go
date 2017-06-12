package informer

// ObserverFuncs implements Observer interface.
type ObserverFuncs struct {
	OnListFunc  func()
	OnWatchFunc func()
}

func (o ObserverFuncs) OnList() {
	o.OnListFunc()
}

func (o ObserverFuncs) OnWatch() {
	o.OnWatchFunc()
}
