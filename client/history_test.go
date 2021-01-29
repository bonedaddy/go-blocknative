package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgHistory(t *testing.T) {
	hist := &MsgHistory{}
	type arg struct {
		num int
	}
	set := func() {
		a1 := arg{1}
		a2 := arg{2}
		a3 := arg{3}
		hist.Push(a1)
		require.Equal(t, 1, hist.Len())
		hist.Push(a2)
		require.Equal(t, 2, hist.Len())
		hist.Push(a3)
		require.Equal(t, 3, hist.Len())
	}

	// test pushes
	set()
	// test pop
	item := hist.Pop().(arg)
	require.Equal(t, 2, hist.Len()) // test length
	require.Equal(t, 1, item.num)   // test value
	item = hist.Pop().(arg)
	require.Equal(t, 1, hist.Len()) // test length
	require.Equal(t, 2, item.num)   // test value
	item = hist.Pop().(arg)
	require.Equal(t, 0, hist.Len()) // test length
	require.Equal(t, 3, item.num)   // test value
	hist.Push(arg{10})
	require.Equal(t, 1, hist.Len())
	item = hist.Pop().(arg)
	require.Equal(t, 0, hist.Len())
	require.Equal(t, 10, item.num)
	require.Equal(t, nil, hist.Pop())
	// set msg history again
	set()
	all := hist.PopAll()
	require.Equal(t, 0, hist.Len())
	for i := 0; i < len(all); i++ {
		item = all[i].(arg)
		switch i {
		case 0:
			require.Equal(t, 1, item.num)
		case 1:
			require.Equal(t, 2, item.num)
		case 2:
			require.Equal(t, 3, item.num)
		}
	}
}
