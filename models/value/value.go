package value

type T uint16

const (
	ValueTypeUnknown        T = 0
	ValueTypeFloatWithScale T = 1
)

func Int(v int64) V {
	return V{
		T:     ValueTypeFloatWithScale,
		V:     v,
		Scale: 1,
	}
}

func IntWithScale(v int64, scale int64) V {
	return V{
		T:     ValueTypeFloatWithScale,
		V:     v,
		Scale: scale,
	}
}

type V struct {
	T     T
	V     int64
	Scale int64
}
