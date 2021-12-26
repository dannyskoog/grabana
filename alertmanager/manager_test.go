package alertmanager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultContactPoint(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
		DefaultContactPoint("team-a"),
	)

	req.Equal("team-a", manager.builder.Config.Route.Receiver)
}

func TestDefaultContactPointCanBeImplicit(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
	)

	req.Equal("team-b", manager.builder.Config.Route.Receiver)
}

func TestContactPoints(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
	)

	req.Len(manager.builder.Config.Receivers, 2)
}

func TestRouting(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-a"),
		),
		Routing(
			Policy("team-a", TagEq("owner", "team-a")),
		),
	)

	req.Len(manager.builder.Config.Route.Routes, 1)
}

func TestMarshalJSON(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-a"),
		),
		Routing(
			Policy("team-a", TagEq("owner", "team-a")),
		),
	)

	_, errJSON := manager.MarshalJSON()
	_, errJSONIndent := manager.MarshalIndentJSON()

	req.NoError(errJSON)
	req.NoError(errJSONIndent)
}
