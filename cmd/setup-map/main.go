package main
/*
main.go invokes a few examples of the concurrent-map supported
operations. For more details, view the cmap directory for supported
operations.
 */

import (
	"fmt"
	"github.com/vishjain/concurrent-map/cmap"
)

func main() {
	csmap := cmap.New(5)
	ok := csmap.SetIfAbsent("a", "b")
	fmt.Printf("Insertion of a was successful %v \n", ok)
	val, _ := csmap.Get("a")
	fmt.Printf("Value of a is %v \n", val.(string))

	// Remove from map according to callback.
	removeCallbackFn := func(key string, v interface{}, valueExists bool) bool {
		return valueExists
	}
	removalStatus := csmap.RemoveCb("a", removeCallbackFn)
	fmt.Printf("Was the removal successful %v \n", removalStatus)
	fmt.Printf("Is concurrent map empty? %v \n", csmap.IsEmpty())
}
