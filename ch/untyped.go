package ch

//go:generate genny -in=generic.genny -out=typed.go gen "ChannelType=BUILTINS"

// MakeAny takes the given value and returns a channel with the value
// already sent to it, ready to be read.
//
// Not type safe, use with care.
func MakeAny(v interface{}) <-chan interface{} {
	channel := make(chan interface{}, 1)
	channel <- v
	return channel
}

// MakeAnySlice takes the given value and returns a channel with the value
// already sent to it, ready to be read.
//
// Not type safe, use with care.
func MakeAnySlice(v []interface{}) <-chan []interface{} {
	channel := make(chan []interface{}, 1)
	channel <- v
	return channel
}
