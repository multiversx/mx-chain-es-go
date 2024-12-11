package runType

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-es-indexer-go/factory"
)

func createComponents() (factory.RunTypeComponentsHandler, error) {
	rtcf := NewRunTypeComponentsFactory()
	return NewManagedRunTypeComponents(rtcf)
}

func TestNewManagedRunTypeComponents(t *testing.T) {
	t.Parallel()

	t.Run("should error", func(t *testing.T) {
		managedRunTypeComponents, err := NewManagedRunTypeComponents(nil)
		require.ErrorIs(t, err, errNilRunTypeComponents)
		require.Nil(t, managedRunTypeComponents)
	})
	t.Run("should work", func(t *testing.T) {
		rtcf := NewRunTypeComponentsFactory()
		managedRunTypeComponents, err := NewManagedRunTypeComponents(rtcf)
		require.NoError(t, err)
		require.False(t, managedRunTypeComponents.IsInterfaceNil())
	})
}

func TestManagedRunTypeComponents_Create(t *testing.T) {
	t.Parallel()

	t.Run("should work with getters", func(t *testing.T) {
		t.Parallel()

		managedRunTypeComponents, err := createComponents()
		require.NoError(t, err)

		err = managedRunTypeComponents.Create()
		require.NoError(t, err)

		require.Equal(t, runTypeComponentsName, managedRunTypeComponents.String())
		require.NoError(t, managedRunTypeComponents.Close())
	})
}

func TestManagedRunTypeComponents_Close(t *testing.T) {
	t.Parallel()

	managedRunTypeComponents, _ := createComponents()
	require.NoError(t, managedRunTypeComponents.Close())

	err := managedRunTypeComponents.Create()
	require.NoError(t, err)

	require.NoError(t, managedRunTypeComponents.Close())
}

func TestManagedRunTypeComponents_CheckSubcomponents(t *testing.T) {
	t.Parallel()

	managedRunTypeComponents, _ := createComponents()
	err := managedRunTypeComponents.CheckSubcomponents()
	require.Equal(t, errNilRunTypeComponents, err)

	err = managedRunTypeComponents.Create()
	require.NoError(t, err)

	//TODO check for nil each subcomponent - MX-15371
	err = managedRunTypeComponents.CheckSubcomponents()
	require.NoError(t, err)

	require.NoError(t, managedRunTypeComponents.Close())
}
