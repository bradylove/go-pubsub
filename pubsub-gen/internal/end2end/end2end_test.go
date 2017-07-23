package end2end_test

import (
	"flag"
	"testing"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	"github.com/apoydence/pubsub"
)

func TestEnd2End(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)
	flag.Parse()

	o.Spec("routes data as expected", func(t *testing.T) {
		ps := pubsub.New()
		traverser := StructTraverser{}
		sub1 := &mockSubscription{}
		sub2 := &mockSubscription{}
		sub3 := &mockSubscription{}
		sub4 := &mockSubscription{}
		ps.Subscribe(sub1, []string{""})
		ps.Subscribe(sub2, []string{"1", "", "y1", "1", "a"})
		ps.Subscribe(sub3, []string{"", "", "y1", "", "b"})
		ps.Subscribe(sub4, []string{"", "", "", "", "", "y2"})

		ps.Publish(&X{i: 1, j: "a", y1: Y{i: 1, j: "a"}, y2: &Y{i: 1, j: "a"}}, traverser)
		ps.Publish(&X{i: 1, j: "a", y1: Y{i: 2, j: "b"}, y2: &Y{i: 1, j: "a"}}, traverser)
		ps.Publish(&X{i: 1, j: "x", y1: Y{i: 2, j: "b"}}, traverser)

		Expect(t, sub1.callCount).To(Equal(3))
		Expect(t, sub2.callCount).To(Equal(1))
		Expect(t, sub3.callCount).To(Equal(2))
		Expect(t, sub3.callCount).To(Equal(2))
	})
}

type mockSubscription struct {
	callCount int
}

func (m *mockSubscription) Write(data interface{}) {
	m.callCount++
}

//go:generate go install github.com/apoydence/pubsub/pubsub-gen
//go:generate $GOPATH/bin/pubsub-gen --struct-name=github.com/apoydence/pubsub/pubsub-gen/internal/end2end.X --package=end2end_test --traverser=StructTraverser --output=$GOPATH/src/github.com/apoydence/pubsub/pubsub-gen/internal/end2end/generated_traverser_test.go --pointer

type X struct {
	i  int
	j  string
	y1 Y
	y2 *Y
}

type Y struct {
	i int
	j string
}
