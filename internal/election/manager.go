package election

import (
	"context"
	"time"

	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

type ElectionManager struct {
	db store.Store
}

func NewElectionManager(db store.Store) *ElectionManager {

	return &ElectionManager{
		db: db,
	}
}

func (em *ElectionManager) StartElectionManagerRoutine() {
	for range time.NewTicker(1 * time.Minute).C {
		em.checkElections()
	}
}

func (em *ElectionManager) checkElections() {
	elections, err := em.db.GetAllElections(context.Background())
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, election := range elections {
		if election.LastRunAt == nil || time.Since(*election.LastRunAt) > time.Duration(election.Hours)*time.Hour {
			log.Infof("stopping any previous election and starting election for channel %s", election.ChannelTwitchID)
			// em.runElection(election)
		}
	}
}
