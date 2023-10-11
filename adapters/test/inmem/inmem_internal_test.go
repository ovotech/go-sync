package inmem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ovotech/go-sync/pkg/types"
)

func TestInMem(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	inMem, err := Init(ctx, map[types.ConfigKey]string{})
	assert.NoError(t, err)

	things, err := inMem.Get(ctx)

	assert.NoError(t, err)
	assert.Equal(t, []string{}, things)

	err = inMem.Add(ctx, []string{"foo", "bar"})
	assert.NoError(t, err)

	things, err = inMem.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"foo", "bar"}, things)

	err = inMem.Remove(ctx, []string{"foo", "buzz"})
	assert.NoError(t, err)

	things, err = inMem.Get(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"bar"}, things)
}
