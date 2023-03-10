package snowid

import (
	"encoding/json"
	"fmt"

	"github.com/xkamail/snowflake"
)

type ID int64

func (i *ID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	id, err := snowflake.ParseString(s)
	if err != nil {
		return err
	}
	*i = ID(id.Int64())
	return nil
}

func (i *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprint(*i))
}

func (i *ID) Int() int64 {
	return int64(*i)
}

const Zero = ID(0)

var node, _ = snowflake.NewNode(1)

func Gen() ID {
	return ID(node.Generate().Int64())
}
