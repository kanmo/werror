package werror

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"runtime"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("werror")
	assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.Unknown), err.Error())
	assert.Equal(t, true, err.(*Error).shouldReport)
	callers := err.(*Error).Callers()
	fn := runtime.FuncForPC(callers[0])
	assert.Equal(t, "werror.TestNew", fn.Name())
}

type testErrorReason struct{}

func (_ testErrorReason) String() string {
	return "test"
}

func TestWrap(t *testing.T) {
	t.Run("wrap werror", func(t *testing.T) {
		werr := wrapTestError(errors.New("werror"))
		assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.Unknown), werr.Error())
		assert.Equal(t, true, werr.(*Error).shouldReport)
		callers := werr.(*Error).Callers()
		fn := runtime.FuncForPC(callers[0])
		assert.Equal(t, "werror.wrapTestError", fn.Name())
	})

	t.Run("wrap nil", func(t *testing.T) {
		werr := Wrap(nil)
		assert.Nil(t, werr)
	})

	t.Run("already wrapped", func(t *testing.T) {
		origCallers := []uintptr{1, 2, 3, 4, 5}
		var werr error = &Error{callers: origCallers}
		werr = Wrap(werr)
		callers := werr.(*Error).Callers()
		assert.Equal(t, origCallers, callers)
	})
}

func TestWithCode(t *testing.T) {
	t.Run("specify code", func(t *testing.T) {
		err := errors.New("werror")
		err = Wrap(err, WithCode(codes.InvalidArgument))
		assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.InvalidArgument), err.Error())
	})

	t.Run("not specify code", func(t *testing.T) {
		err := errors.New("werror")
		err = Wrap(err)
		assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.Unknown), err.Error())
	})
}

func TestWithReason(t *testing.T) {
	t.Run("specify reason", func(t *testing.T) {
		err := errors.New("werror")
		err = Wrap(err, WithReason(testErrorReason{}))
		assert.Equal(t, fmt.Sprintf("reason: test code: %d message: werror", codes.Unknown), err.Error())
	})

	t.Run("not specify reason", func(t *testing.T) {
		err := errors.New("werror")
		err = Wrap(err)
		assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.Unknown), err.Error())
	})

	t.Run("specifiy invalid type reason", func(t *testing.T) {
		err := errors.New("werror")
		err = Wrap(err, WithReason("test"))
		assert.Equal(t, fmt.Sprintf("reason:  code: %d message: werror", codes.Unknown), err.Error())
	})
}

func wrapTestError(err error) error {
	return Wrap(err)
}
