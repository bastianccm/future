package future

import "context"

func Await[T any](ctx context.Context, future func(ctx context.Context) (T, error)) func() (T, error) {
	valC := make(chan T)
	errC := make(chan error)

	go func() {
		val, err := future(ctx)
		if err != nil {
			errC <- err
			return
		}
		valC <- val
	}()

	return func() (T, error) {
		select {
		case val := <-valC:
			return val, nil
		case err := <-errC:
			return *new(T), err
		case <-ctx.Done():
			return *new(T), ctx.Err()
		}
	}
}

func Resolve1[T any](cb1 func() (T, error)) (T, error) {
	return cb1()
}

func Resolve2[T, U any](cb1 func() (T, error), cb2 func() (U, error)) (T, U, error) {
	v1, err := cb1()
	if err != nil {
		return v1, *new(U), err
	}
	v2, err := Resolve1(cb2)
	return v1, v2, err
}

func Resolve3[T, U, V any](cb1 func() (T, error), cb2 func() (U, error), cb3 func() (V, error)) (T, U, V, error) {
	v1, err := cb1()
	if err != nil {
		return v1, *new(U), *new(V), err
	}
	v2, v3, err := Resolve2(cb2, cb3)
	return v1, v2, v3, err
}

func Resolve4[T, U, V, W any](cb1 func() (T, error), cb2 func() (U, error), cb3 func() (V, error), cb4 func() (W, error)) (T, U, V, W, error) {
	v1, err := cb1()
	if err != nil {
		return v1, *new(U), *new(V), *new(W), err
	}
	v2, v3, v4, err := Resolve3(cb2, cb3, cb4)
	return v1, v2, v3, v4, err
}

func Resolve5[T, U, V, W, X any](cb1 func() (T, error), cb2 func() (U, error), cb3 func() (V, error), cb4 func() (W, error), cb5 func() (X, error)) (T, U, V, W, X, error) {
	v1, err := cb1()
	if err != nil {
		return v1, *new(U), *new(V), *new(W), *new(X), err
	}
	v2, v3, v4, v5, err := Resolve4(cb2, cb3, cb4, cb5)
	return v1, v2, v3, v4, v5, err
}

func Resolve6[T, U, V, W, X, Y any](cb1 func() (T, error), cb2 func() (U, error), cb3 func() (V, error), cb4 func() (W, error), cb5 func() (X, error), cb6 func() (Y, error)) (T, U, V, W, X, Y, error) {
	v1, err := cb1()
	if err != nil {
		return v1, *new(U), *new(V), *new(W), *new(X), *new(Y), err
	}
	v2, v3, v4, v5, v6, err := Resolve5(cb2, cb3, cb4, cb5, cb6)
	return v1, v2, v3, v4, v5, v6, err
}
