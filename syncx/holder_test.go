package syncx

import (
	"fmt"
	"testing"
)

func TestHolder(t *testing.T) {
	var holder Holder[int]
	var g Group
	holder.Set(1)
	g.Go(func() error {
		holder.Set(2)
		fmt.Println(
			holder.Get(),
		)
		return nil
	})
	g.Go(func() error {
		fmt.Println(
			holder.Get(),
		)
		return nil
	})

	g.Wait()
	fmt.Println(holder.Get())

}
