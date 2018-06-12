package blc

type Walets struct {

	Walets map[string]*Walet
}

func NewWalets() *Walets {
	walets := &Walets{}
	walets.Walets = make(map[string]*Walet)
	return walets
}

func (w *Walets)CreatWalets() {
	walet := NewWalet()

	w.Walets[string(walet.GetAddress())] = walet

}