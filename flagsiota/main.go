package main

type Flag uint

const (
	FlagUp Flag = 1 << iota
	FlagBroadcast
	FlagLoopback
	FlagPointToPoint
	FlagMulticast
)

func main() {

}

func IsUp(v Flag) bool     { return v&FlagUp == FlagUp }
func TurnDown(v *Flag)     { *v &^= FlagUp }
func SetBroadcast(v *Flag) { *v |= FlagBroadcast }
func IsCast(v Flag) bool {
	// return v & (FlagBroadcast | )
	return true
}
