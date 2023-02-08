package device

const (
	Deng       = iota
	Kaiguan    = iota
	Chuanganqi = iota
	Men        = iota
)

const (
	Switch = iota
	Slider = iota
	Text   = iota
)

type Entry struct {
	Type  uint
	Value string
}

type Resp struct {
	ID     uint
	PosX   float64
	PosY   float64
	Zoom   float64
	Img    string
	Name   string
	Status interface{}
	Type   uint
	Sid    uint
}
