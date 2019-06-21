package goro

// Backend is an interface through which a Screen is displayed and controlled.
type Backend interface {
	Init() error
	Run(func(*Screen)) error
	Refresh()
	Quit()
	Size() (int, int)
	SetSize(int, int)
	Scale() float64
	SetScale(float64)
	SetTitle(string)
}
