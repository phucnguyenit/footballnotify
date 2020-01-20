package types

import "fmt"

// Event ...
type Event struct {
	MatchID                    string       `json:"match_id"`
	CountryID                  string       `json:"country_id"`
	CountryName                string       `json:"country_name"`
	LeagueID                   string       `json:"league_id"`
	LeagueName                 string       `json:"league_name"`
	MatchStatus                string       `json:"match_status"`
	MatchDate                  string       `json:"match_date"`
	MatchTime                  string       `json:"match_time"`
	MatchHomeTeamName          string       `json:"match_hometeam_name"`
	MatchHomeTeamScore         string       `json:"match_hometeam_score"`
	MatchAwayTeamName          string       `json:"match_awayteam_name"`
	MatchAwayTeamScore         string       `json:"match_awayteam_score"`
	MatchHomeTeamHaltTimeScore string       `json:"match_hometeam_halttime_score"`
	MatchAwayTeamHaltTimeScore string       `json:"match_awayteam_halttime_score"`
	MatchHomeTeamExtraScore    string       `json:"match_hometeam_extra_score"`
	MatchAwayTeamExtraScore    string       `json:"match_awayteam_extra_score"`
	MatchHomeTeamSystem        string       `json:"match_hometeam_system"`
	MatchAwayTeamSystem        string       `json:"match_awayteam_system"`
	MatchLive                  string       `json:"match_live"`
	GoalScorer                 []GoalScorer `json:"goalscorer"`
	Substitutions              struct {
		Home []LineUp `json:"home"`
		Away []LineUp `json:"away"`
	} `json:"substitutions"`
	LineUp struct {
		Home struct {
			StartingLineUps []LineUp `json:"starting_lineups"`
			Substitutes     []LineUp `json:"substitutes"`
			Coach           []LineUp `json:"coach"`
		} `json:"home"`
		Away struct {
			StartingLineUps []LineUp `json:"starting_lineups"`
			Substitutes     []LineUp `json:"substitutes"`
			Coach           []LineUp `json:"coach"`
		} `json:"away"`
	} `json:"lineup"`
}

// IsLive ...
func (e Event) IsLive(ne Event) bool {
	return e.MatchLive == "0" && ne.MatchLive == "1"
}

// GoalScorerChanges ...
func (e Event) GoalScorerChanges(ne Event) []GoalScorer {
	if len(e.GoalScorer) < len(ne.GoalScorer) {
		return ne.GoalScorer[len(e.GoalScorer)+1:]
	}
	return nil
}

// HomeSubChanges ...
func (e Event) HomeSubChanges(ne Event) []LineUp {
	if len(e.Substitutions.Home) < len(ne.Substitutions.Home) {
		return ne.Substitutions.Home[len(e.Substitutions.Home):]
	}
	return nil
}

// AwayTeamSubChanges ...
func (e Event) AwayTeamSubChanges(ne Event) []LineUp {
	if len(e.Substitutions.Away) < len(ne.Substitutions.Away) {
		return ne.Substitutions.Away[len(e.Substitutions.Away):]
	}
	return nil
}

// IsEnd ...
func (e Event) IsEnd(ne Event) bool {
	return e.MatchStatus != "FT" && ne.MatchStatus == "FT"
}

// GetNotificationMessages ...
func (e Event) GetNotificationMessages(ne Event) []Message {
	topic := "team:" + e.MatchID
	msgs := []Message{}

	if e.IsLive(ne) == true {
		msgs = append(msgs, Message{
			Topic: topic,
			Title: fmt.Sprintf("Trận đấu giữa %s và %s đã bắt đầu",
				e.MatchHomeTeamName, e.MatchAwayTeamName),
		})
	}

	if e.IsEnd(ne) == true {
		msgs = append(msgs, Message{
			Topic: topic,
			Title: fmt.Sprintf("Trận đấu giữa %s và %s đã kết thúc với tỷ số (%s-%s)",
				ne.MatchHomeTeamName, ne.MatchAwayTeamName, ne.MatchHomeTeamScore, ne.MatchAwayTeamScore),
		})
	}

	goalscorerChanges := e.GoalScorerChanges(ne)
	if len(goalscorerChanges) > 0 {
		for _, goalscorerChange := range goalscorerChanges {
			msgs = append(msgs, Message{
				Topic: topic,
				Title: fmt.Sprintf("Trận đấu (%s-%s) %s đã ghi bàn tỷ số (%s-%s) ",
					ne.MatchHomeTeamName,
					ne.MatchAwayTeamName,
					goalscorerChange.GetScorerName(),
					ne.MatchHomeTeamScore,
					ne.MatchAwayTeamScore),
			})
		}
	}

	homeSubChanges := e.HomeSubChanges(ne)
	if len(homeSubChanges) > 0 {
		for _, subChange := range homeSubChanges {
			msgs = append(msgs, Message{
				Topic: topic,
				Title: fmt.Sprintf("Trận đấu (%s-%s) Thay người %s bên đội %s",
					ne.MatchHomeTeamName,
					ne.MatchAwayTeamName,
					subChange.LineUpPlayer,
					ne.MatchHomeTeamName),
			})
		}
	}

	awaySubChanges := e.AwayTeamSubChanges(ne)
	if len(awaySubChanges) > 0 {
		for _, subChange := range awaySubChanges {
			msgs = append(msgs, Message{
				Topic: topic,
				Title: fmt.Sprintf("Trận đấu (%s-%s) Thay người %s bên đội %s",
					ne.MatchHomeTeamName,
					ne.MatchAwayTeamName,
					subChange.LineUpPlayer,
					ne.MatchAwayTeamName),
			})
		}
	}

	return msgs
}

// GoalScorer ...
type GoalScorer struct {
	Time       string `json:"time"`
	HomeScorer string `json:"home_scorer"`
	Score      string `json:"score"`
	AwayScorer string `json:"away_scorer"`
}

// GetScorerName ...
func (g GoalScorer) GetScorerName() string {
	if g.HomeScorer != "" {
		return g.HomeScorer
	}
	return g.AwayScorer
}

// Card ...
type Card struct {
	Time      string `json:"time"`
	HomeFault string `json:"home_fault"`
	Score     string `json:"score"`
	AwayFault string `json:"away_fault"`
}

// LineUp ...
type LineUp struct {
	LineUpPlayer   string `json:"lineup_player"`
	LineUpNumber   string `json:"lineup_number"`
	LineUpPosition string `json:"lineup_position"`
	LineTime       string `json:"line_time"`
}

// Events ...
type Events []Event

// GetNotificationMessages ...
func (events Events) GetNotificationMessages(nEvents Events) []Message {
	msgs := []Message{}
	for idx, event := range events {
		msgs = append(msgs, event.GetNotificationMessages(nEvents[idx])...)
	}
	return msgs
}
