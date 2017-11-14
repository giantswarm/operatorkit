package tprtest

type TPRTest struct{}

func New() *TPRTest {
	return &TPRTest{}
}

func (i *TPRTest) CreateAndWait() error {
	return nil
}
