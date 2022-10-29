package media

import (
	"testing"

	"github.com/gempir/gempbot/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestCanRegisterConnectionAndHandleJoin(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore())

	connId := mgr.RegisterConnection("conn1", func(message []byte) {})
	mgr.HandleJoin(connId, "userId1", "")

	assert.Equal(t, 1, mgr.connections.Size())

	val, ok := mgr.rooms.Load("userId1")
	assert.True(t, ok)

	_, ok = val.users.Load(connId)
	assert.True(t, ok)
}

func TestAbortsJoinWhenNoConnectionFound(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore())

	mgr.HandleJoin("conn1", "userId1", "")

	assert.Equal(t, 0, mgr.connections.Size())
	assert.Equal(t, 0, mgr.rooms.Size())
}

func TestCanCreateRoom(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore())

	_ = mgr.getRoom("userId1")
	assert.Equal(t, 1, mgr.rooms.Size())
}

func TestCanGetExistingRoom(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore())

	room := mgr.getRoom("userId1")
	room.CurrentTime = 10
	assert.Equal(t, 1, mgr.rooms.Size())

	room = mgr.getRoom("userId1")
	assert.Equal(t, float32(10), room.CurrentTime)
}

func TestCanHandleTimeChange(t *testing.T) {
	mgr := NewMediaManager(store.NewMockStore())

	connId := mgr.RegisterConnection("conn1", func(message []byte) {})
	mgr.HandleJoin(connId, "userId1", "")

	mgr.HandleTimeChange("conn1", "userId1", "videoId1", 10)
	room := mgr.getRoom("userId1")
	assert.Equal(t, float32(10), room.CurrentTime)
}
