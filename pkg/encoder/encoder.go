package encoder

import "io"

type Encoder interface {
	Encode(interface{}, io.Writer) error
}
