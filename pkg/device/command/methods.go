package command

type Method string

const (
	GET_PROP   Method = "get_prop"
	SET_RGB    Method = "set_rgb"
	SET_HSV    Method = "set_hsv"
	SET_BRIGHT Method = "set_bright"
	SET_POWER  Method = "set_power"
	TOGGLE     Method = "toggle"
)
