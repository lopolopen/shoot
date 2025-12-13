package enums

//go:generate go tool shoot enum -json -text -file=$GOFILE

type OrderStatus int32

const (
	OrderStatusPending OrderStatus = iota
	OrderStatusCompleted
	OrderStatusCanceled
)
