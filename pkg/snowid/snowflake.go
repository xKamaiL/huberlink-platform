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

func (i *ID) String() string {
	return fmt.Sprint(*i)
}

const Zero = ID(0)

var node, _ = snowflake.NewNode(1)

func Gen() ID {
	return ID(node.Generate().Int64())
}

func Parse(str string) (ID, error) {
	id, err := snowflake.ParseString(str)
	if err != nil {
		return Zero, err
	}
	return ID(id.Int64()), nil
}
