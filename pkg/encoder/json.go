package encoder

import (
	"encoding/json"
	"io"
)

type JsonEncoder struct{}

func (e JsonEncoder) Encode(v interface{}, w io.Writer) error {
	if b, err := json.Marshal(v); err != nil {
		return err
	} else {
		w.Write(b)
		return nil
	}
}
