package main

import (
	"log"
	"strings"

	"github.com/apoydence/pubsub"
)

//         a     x
//         |      \
//         b       y
//        / \       \
//       c   d       z
func main() {
	ps := pubsub.New()

	ps.Subscribe(Subscription("sub-1"), []string{"a", "b", "c"})
	ps.Subscribe(Subscription("sub-2"), []string{"a", "b", "d"})
	ps.Subscribe(Subscription("sub-3"), []string{"a", "b", "e"})
	ps.Subscribe(Subscription("sub-4"), []string{"a"})
	ps.Subscribe(Subscription("sub-5"), []string{"a", "b"})
	ps.Subscribe(Subscription("sub-6"), []string{"x", "y", "z"})

	dataMap1 := map[string][]string{
		"":      []string{"a"},
		"a":     []string{"b"},
		"a-b":   []string{"c", "d"},
		"a-b-c": nil,
		"a-b-d": nil,
	}
	ps.Publish("data-1", StaticAssigner(dataMap1))

	dataMap2 := map[string][]string{
		"":      []string{"x"},
		"x":     []string{"y"},
		"x-y":   []string{"z"},
		"x-y-z": nil,
	}
	ps.Publish("data-2", StaticAssigner(dataMap2))
}

// Subscription writes any results to stderr
type Subscription string

// Write implements pubsub.Subscription
func (s Subscription) Write(data interface{}) {
	log.Printf("%s <- %s", s, data)
}

// StaticAssigner assigns data based on its underlying map and not the data.
// Therefore, it does not look at the data to decide where the data belongs.
// Only the given path.
type StaticAssigner map[string][]string

func (a StaticAssigner) Assign(data interface{}, currentPath []string) pubsub.AssignedPaths {
	path := strings.Join(currentPath, "-")
	ps, ok := a[path]
	if !ok {
		log.Panicf("Unknown path: '%s'", path)
	}

	return pubsub.Paths(ps)
}
