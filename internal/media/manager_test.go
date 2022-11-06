package media

import (
	"testing"

	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestCanRegisterConnectionAndHandleJoin(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore(), helixclient.NewMockClient(), bot.NewMockbot())

	connId := mgr.RegisterConnection("conn1", func(message []byte) {})
	mgr.HandleJoin(connId, "userId1", "")

	assert.Equal(t, 1, mgr.connections.Size())

	val, ok := mgr.rooms.Load("userId1")
	assert.True(t, ok)

	_, ok = val.users.Load(connId)
	assert.True(t, ok)
}

func TestAbortsJoinWhenNoConnectionFound(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore(), helixclient.NewMockClient(), bot.NewMockbot())

	mgr.HandleJoin("conn1", "userId1", "channel")

	assert.Equal(t, 0, mgr.connections.Size())
	assert.Equal(t, 0, mgr.rooms.Size())
}

func TestCanCreateRoom(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore(), helixclient.NewMockClient(), bot.NewMockbot())

	_ = mgr.getRoom("userId1")
	assert.Equal(t, 1, mgr.rooms.Size())
}

func TestCanGetExistingRoom(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore(), helixclient.NewMockClient(), bot.NewMockbot())

	room := mgr.getRoom("userId1")
	room.Time = 10
	assert.Equal(t, 1, mgr.rooms.Size())

	room = mgr.getRoom("userId1")
	assert.Equal(t, float32(10), room.Time)
}

func TestCanHandlePlayerStateChange(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore(), helixclient.NewMockClient(), bot.NewMockbot())

	connId := mgr.RegisterConnection("conn1", func(message []byte) {})
	mgr.HandleJoin(connId, "userId1", "")

	mgr.HandlePlayerState("conn1", "userId1", PLAYING, "videoId1", 10)
	room := mgr.getRoom("userId1")
	assert.Equal(t, float32(10), room.Time)
}
