package future_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bastianccm/future"
	"github.com/stretchr/testify/assert"
)

func TestFutureTyped(t *testing.T) {
	ctx := context.Background()

	a, b, c, d, e, f, err := future.Resolve6(
		future.Promise(ctx, func(ctx context.Context) (int8, error) { return 1, nil }),
		future.Promise(ctx, func(ctx context.Context) (uint8, error) { return 1, nil }),
		future.Promise(ctx, func(ctx context.Context) (int16, error) { return 1, nil }),
		future.Promise(ctx, func(ctx context.Context) (uint16, error) { return 1, nil }),
		future.Promise(ctx, func(ctx context.Context) (int32, error) { return 1, nil }),
		future.Promise(ctx, func(ctx context.Context) (uint32, error) { return 1, nil }),
	)

	assert.NoError(t, err)
	assert.Equal(t, int8(1), a)
	assert.Equal(t, uint8(1), b)
	assert.Equal(t, int16(1), c)
	assert.Equal(t, uint16(1), d)
	assert.Equal(t, int32(1), e)
	assert.Equal(t, uint32(1), f)

	x, err := future.ResolveN(
		future.Promise(ctx, func(context.Context) (int, error) { return 1, nil }),
		future.Promise(ctx, func(context.Context) (int, error) { return 2, nil }),
		future.Promise(ctx, func(context.Context) (int, error) { return 3, nil }),
		future.Promise(ctx, func(context.Context) (int, error) { return 4, nil }),
		future.Promise(ctx, func(context.Context) (int, error) { return 5, nil }),
	)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, x)

	x, err = future.ResolveN(
		future.Promise(ctx, func(context.Context) (int, error) { return 1, nil }),
		future.Promise(ctx, func(context.Context) (int, error) { return 2, errors.New("err") }),
	)
	assert.Error(t, err)
	assert.Nil(t, x)
}

func TestFutureErrors(t *testing.T) {
	ctx := context.Background()
	errTest := errors.New("test error")

	cases := []struct {
		a, b, c, d, e, f                   int
		errA, errB, errC, errD, errE, errF error
		ea, eb, ec, ed, ee, ef             int
	}{
		{1, 2, 3, 4, 5, 6, nil, nil, nil, nil, nil, errTest, 1, 2, 3, 4, 5, 0},
		{1, 2, 3, 4, 5, 6, nil, nil, nil, nil, errTest, nil, 1, 2, 3, 4, 0, 0},
		{1, 2, 3, 4, 5, 6, nil, nil, nil, errTest, nil, nil, 1, 2, 3, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, nil, nil, errTest, nil, nil, nil, 1, 2, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, nil, errTest, nil, nil, nil, nil, 1, 0, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, errTest, nil, nil, nil, nil, nil, 0, 0, 0, 0, 0, 0},
	}
	for _, tc := range cases {
		a, b, c, d, e, f, err := future.Resolve6(
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.a, tc.errA }),
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.b, tc.errB }),
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.c, tc.errC }),
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.d, tc.errD }),
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.e, tc.errE }),
			future.Promise(ctx, func(ctx context.Context) (int, error) { return tc.f, tc.errF }),
		)

		assert.ErrorIs(t, err, errTest)
		assert.Equal(t, tc.ea, a)
		assert.Equal(t, tc.eb, b)
		assert.Equal(t, tc.ec, c)
		assert.Equal(t, tc.ed, d)
		assert.Equal(t, tc.ee, e)
		assert.Equal(t, tc.ef, f)
	}
}

func TestPromiseCancelled(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err := future.Promise(ctx, func(ctx context.Context) (int, error) { return 0, nil })()
	assert.ErrorIs(t, err, context.Canceled)
}
