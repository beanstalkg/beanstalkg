package architecture

import (
	"fmt"
	"testing"
)

func TestCommand_Parse(t *testing.T) {
	c := Command{}
	c.Parse("use test")
	fmt.Println(c)
	if c.RawCommand != "use test" {
		fmt.Println("First")
		t.Fail()
	}
	if c.Name != "use" {
		fmt.Println("Second")
		t.Fail()
	}

	c = Command{}
	done, err := c.Parse("put 1 2 3 4")

	if err != nil {
		fmt.Println("Third")
		t.Fail()
	}

	if done {
		fmt.Println("Fourth")
		t.Fail()
	}

	if c.Name != "put" || c.Params["ttr"] != "3" {
		fmt.Println("Fifth")
		t.Fail()
	}

	c.Parse("data")

	if c.Params["data"] != "data" {
		fmt.Println("Sixth")
		t.Fail()
	}

	c = Command{}
	kicked, err := c.Parse("kick 4")
	if err != nil {
		t.Fail()
	}
	if !kicked {
		fmt.Println("failed to kicked")
		t.Fail()
	}
}
