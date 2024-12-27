package main

type Team struct {
	TeamID string `json:"team_id"`
	Rounds struct {
		Won  int `json:"won"`
		Lost int `json:"lost"`
	}
	Won bool `json:"won"`
}

type Player struct {
	Name  string `json:"name"`
	Tag   string `json:"tag"`
	Team  string `json:"team_id"`
	Agent struct {
		Name string `json:"name"`
	} `json:"agent"`
	Stats struct {
		Score   int `json:"score"`
		Kills   int `json:"kills"`
		Deaths  int `json:"deaths"`
		Assists int `json:"assists"`
	} `json:"stats"`
	Tier struct {
		Name string `json:"name"`
	} `json:"tier"`
}

type MatchDataResponse struct {
	Status    int         `json:"status"`
	MatchData []MatchData `json:"data"`
}

type MMRData struct {
	CurrentRR int
	RRChange  int
	Tier      string
}

type MMRDataResponse struct {
	Status int `json:"status"`
	Data   struct {
		Current struct {
			Tier struct {
				Name string `json:"name"`
			} `json:"tier"`
			RR         int `json:"rr"`
			LastChange int `json:"last_change"`
		} `json:"current"`
	} `json:"data"`
}

type WebhookData struct {
	Embeds []Embed `json:"embeds"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Color       int          `json:"color"`
	Timestamp   string       `json:"timestamp"`
	Image       EmbedImage   `json:"image"`
	Fields      []EmbedField `json:"fields"`
	Footer      EmbedFooter  `json:"footer"`
}

type EmbedImage struct {
	URL   string `json:"url"`
	Width int    `json:"height"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type MatchData struct {
	Metadata struct {
		MatchID string `json:"match_id"`
		Map     struct {
			Name string `json:"name"`
		} `json:"map"`
		StartedAt string `json:"started_at"`
	} `json:"metadata"`
	Players []Player
	Teams   []Team `json:"teams"`
}

type TrackedPlayerData struct {
	Name     string
	Tag      string
	Platform string
	Region   string
}
