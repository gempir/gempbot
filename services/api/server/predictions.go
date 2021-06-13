package server

import (
	"time"

	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
)

func (s *Server) subscribePredictions(userID string) {
	// response, err := s.helixClient.CreateEventSubSubscription(userID, s.cfg.WebhookApiBaseUrl+"/api/prediction/begin", "channel.prediction.begin")
	// if err != nil {
	// 	log.Errorf("Error subscribing: %s", err)
	// 	return
	// }

	// log.Infof("[%d] created subscription %s", response.StatusCode, response.Error)
	// for _, sub := range response.Data.EventSubSubscriptions {
	// 	s.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	// }

	// response, err = s.helixClient.CreateEventSubSubscription(userID, s.cfg.WebhookApiBaseUrl+"/api/prediction/lock", "channel.prediction.lock")
	// if err != nil {
	// 	log.Errorf("Error subscribing: %s", err)
	// 	return
	// }

	// log.Infof("[%d] created subscription %s", response.StatusCode, response.Error)
	// for _, sub := range response.Data.EventSubSubscriptions {
	// 	s.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	// }

	response, err := s.helixClient.CreateEventSubSubscription(userID, s.cfg.WebhookApiBaseUrl+"/api/prediction/end", "channel.prediction.end")
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	log.Infof("[%d] created subscription %s", response.StatusCode, response.Error)
	for _, sub := range response.Data.EventSubSubscriptions {
		s.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type)
	}
}

type predictionBegin struct {
	Subscription struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Status    string `json:"status"`
		Cost      int    `json:"cost"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"subscription"`
	Event struct {
		ID                   string `json:"id"`
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		Title                string `json:"title"`
		Outcomes             []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Color string `json:"color"`
		} `json:"outcomes"`
		StartedAt time.Time `json:"started_at"`
		LocksAt   time.Time `json:"locks_at"`
	} `json:"event"`
}

func (s *Server) handlePredictionBegin(c echo.Context) error {
	var data predictionBegin
	done, err := s.handleWebhook(c, &data)
	if err != nil || done {
		return err
	}

	log.Infof("predictionBegin %s", data.Event.LocksAt)
	if data.Event.ID == "" {
		return nil
	}

	err = s.db.SavePrediction(store.PredictionLog{ID: data.Event.ID, OwnerTwitchID: data.Event.BroadcasterUserID, Title: data.Event.Title, StartedAt: data.Event.StartedAt, LockedAt: data.Event.LocksAt})
	if err != nil {
		log.Error(err)
	}

	for _, outcome := range data.Event.Outcomes {
		err = s.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.Event.ID, Title: outcome.Title, Color: outcome.Color})
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

type predictionLock struct {
	Subscription struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Status    string `json:"status"`
		Cost      int    `json:"cost"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"subscription"`
	Event struct {
		ID                   string `json:"id"`
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		Title                string `json:"title"`
		Outcomes             []struct {
			ID            string `json:"id"`
			Title         string `json:"title"`
			Color         string `json:"color"`
			Users         int    `json:"users,omitempty"`
			ChannelPoints int    `json:"channel_points,omitempty"`
			TopPredictors []struct {
				UserName          string      `json:"user_name"`
				UserLogin         string      `json:"user_login"`
				UserID            string      `json:"user_id"`
				ChannelPointsWon  interface{} `json:"channel_points_won"`
				ChannelPointsUsed int         `json:"channel_points_used"`
			} `json:"top_predictors"`
		} `json:"outcomes"`
		StartedAt time.Time `json:"started_at"`
		LockedAt  time.Time `json:"locked_at"`
	} `json:"event"`
}

func (s *Server) handlePredictionLock(c echo.Context) error {
	var data predictionLock
	done, err := s.handleWebhook(c, &data)
	if err != nil || done {
		return err
	}

	log.Infof("predictionLock %s", data.Event.LockedAt)
	if data.Event.ID == "" {
		return nil
	}

	err = s.db.SavePrediction(store.PredictionLog{ID: data.Event.ID, OwnerTwitchID: data.Event.BroadcasterUserID, Title: data.Event.Title, StartedAt: data.Event.StartedAt, LockedAt: data.Event.LockedAt})
	if err != nil {
		log.Error(err)
	}

	for _, outcome := range data.Event.Outcomes {
		err = s.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.Event.ID, Title: outcome.Title, Color: outcome.Color, Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}

type predictionEnd struct {
	Subscription struct {
		ID        string `json:"id"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Status    string `json:"status"`
		Cost      int    `json:"cost"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id"`
		} `json:"condition"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"subscription"`
	Event struct {
		ID                   string `json:"id"`
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		Title                string `json:"title"`
		WinningOutcomeID     string `json:"winning_outcome_id"`
		Outcomes             []struct {
			ID            string `json:"id"`
			Title         string `json:"title"`
			Color         string `json:"color"`
			Users         int    `json:"users"`
			ChannelPoints int    `json:"channel_points"`
			TopPredictors []struct {
				UserName          string `json:"user_name"`
				UserLogin         string `json:"user_login"`
				UserID            string `json:"user_id"`
				ChannelPointsWon  int    `json:"channel_points_won"`
				ChannelPointsUsed int    `json:"channel_points_used"`
			} `json:"top_predictors"`
		} `json:"outcomes"`
		Status    string    `json:"status"`
		StartedAt time.Time `json:"started_at"`
		EndedAt   time.Time `json:"ended_at"`
	} `json:"event"`
}

func (s *Server) handlePredictionEnd(c echo.Context) error {
	var data predictionEnd
	done, err := s.handleWebhook(c, &data)
	if err != nil || done {
		return err
	}

	log.Infof("predictionEnd %s", data.Event.Status)
	if data.Event.ID == "" {
		return nil
	}

	err = s.db.SavePrediction(store.PredictionLog{ID: data.Event.ID, OwnerTwitchID: data.Event.BroadcasterUserID, Title: data.Event.Title, StartedAt: data.Event.StartedAt, EndedAt: data.Event.EndedAt, WinningOutcomeID: data.Event.WinningOutcomeID, Status: data.Event.Status})
	if err != nil {
		log.Error(err)
	}

	for _, outcome := range data.Event.Outcomes {
		err = s.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: data.Event.ID, Title: outcome.Title, Color: outcome.Color, Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
