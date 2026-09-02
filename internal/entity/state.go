package entity

type State string

const (
	StateLayoutType State = "layout_type"
	StatePurpose    State = "purpose"
	StateScale      State = "scale"
	StateSize       State = "size"
	StateMaterial   State = "material"
	StateDetails    State = "details"
	State3DPrint    State = "3d_print"
	StateLandscape  State = "landscape"
	StateDrawings   State = "drawings"
	StateDeadline   State = "deadline"
	StateDelivery   State = "delivery"
	StateAccept     State = "accept"
)
