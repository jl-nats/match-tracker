package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var maps = map[string]string{"Sunset": "https://static.wikia.nocookie.net/valorant/images/5/5c/Loading_Screen_Sunset.png",
	"Lotus":    "https://static.wikia.nocookie.net/valorant/images/d/d0/Loading_Screen_Lotus.png",
	"Pearl":    "https://static.wikia.nocookie.net/valorant/images/a/af/Loading_Screen_Pearl.png",
	"Fracture": "https://static.wikia.nocookie.net/valorant/images/f/fc/Loading_Screen_Fracture.png",
	"Breeze":   "https://static.wikia.nocookie.net/valorant/images/1/10/Loading_Screen_Breeze.png",
	"Icebox":   "https://static.wikia.nocookie.net/valorant/images/1/13/Loading_Screen_Icebox.png",
	"Bind":     "https://static.wikia.nocookie.net/valorant/images/2/23/Loading_Screen_Bind.png",
	"Haven":    "https://static.wikia.nocookie.net/valorant/images/7/70/Loading_Screen_Haven.png",
	"Split":    "https://static.wikia.nocookie.net/valorant/images/d/d6/Loading_Screen_Split.png",
	"Ascent":   "https://static.wikia.nocookie.net/valorant/images/e/e7/Loading_Screen_Ascent.png",
	"Abyss":    "https://www.vpesports.com/wp-content/uploads/2024/06/Screenshot_4-19.png"}

var ranks = map[string]string{"Iron1": "1264549482985558088", "Iron2": "1264549499263516692", "Iron3": "1264549509107552316",
	"Bronze1": "1264549523188092938", "Bronze2": "1264549533682237555", "Bronze3": "1264549543236735076",
	"Silver1": "1264549558583558237", "Silver2": "1264549575205847161", "Silver3": "1264549586333204561",
	"Gold1": "1264549620369850368", "Gold2": "1264549633158418504", "Gold3": "1264549643375612016",
	"Platinum1": "1264549655036039241", "Platinum2": "1264549665429258314", "Platinum3": "1264549676028399677",
	"Diamond1": "1264549711839232031", "Diamond2": "1264549732379004950", "Diamond3": "1264549748329807872",
	"Ascendant1": "1264549759017029743", "Ascendant2": "1264549768433107045", "Ascendant3": "1264549778084073514",
	"Immortal1": "1264549813173616700", "Immortal2": "1264549825878163517", "Immortal3": "1264549838956003412",
	"Radiant": "1264549852940075078", "Unrated": "1264599232631799852"}

func stripSpaces(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func getRankEmoji(tier string) string {
	return "<:" + stripSpaces(tier) + ":" + ranks[stripSpaces(tier)] + ">"
}

var WIN_COLOR = 0x2eb387
var LOSS_COLOR = 0xff2832
var DRAW_COLOR = 0x8BACB5

func CreateEmbedFields(team []Player, label string, totalRounds int) []EmbedField {
	teamFields := []EmbedField{}
	teamFields = append(teamFields, EmbedField{
		Name:   strings.Title(label) + " Team",
		Value:  "",
		Inline: false})
	if totalRounds == 0 {
		totalRounds = 1
	}
	for _, player := range team {
		teamFields = append(teamFields, EmbedField{
			Name:   player.Name + "#```" + player.Tag + "```" + getRankEmoji(player.Tier.Name),
			Value:  "```" + strconv.Itoa(player.Stats.Score/(totalRounds)) + "``` ```" + strconv.Itoa(player.Stats.Kills) + "-" + strconv.Itoa(player.Stats.Deaths) + "-" + strconv.Itoa(player.Stats.Assists) + "``````" + player.Agent.Name + "```",
			Inline: true,
		})
	}
	return teamFields
}

func SeparateTeams(players []Player) ([]Player, []Player) {
	redTeam := []Player{}
	blueTeam := []Player{}

	for _, player := range players {
		if player.Team == "Red" {
			redTeam = append(redTeam, player)
		} else {
			blueTeam = append(blueTeam, player)
		}
	}

	return redTeam, blueTeam
}

func LessFunc(team []Player) func(i, j int) bool {
	return func(i, j int) bool {
		return team[i].Stats.Score > team[j].Stats.Score
	}
}

func CreateEmbed(matchData MatchData, trackedPlayerData TrackedPlayerData, MMRData MMRData) Embed {
	var trackedPlayer Player
	for _, player := range matchData.Players {
		if player.Name == trackedPlayerData.Name && player.Tag == trackedPlayerData.Tag {
			trackedPlayer = player
		}
	}

	embedColor, roundsWon, roundsLost, gameOutcome := processMatchData(matchData, trackedPlayer)

	redTeam, blueTeam := SeparateTeams(matchData.Players)

	sort.Slice(redTeam, LessFunc(redTeam))

	sort.Slice(blueTeam, LessFunc(blueTeam))

	embedFields := append(CreateEmbedFields(redTeam, "red", roundsWon+roundsLost), CreateEmbedFields(blueTeam, "blue", roundsWon+roundsLost)...)

	embed := Embed{
		Title:       "🚨   NEW GAME " + trackedPlayerData.Name + "#" + trackedPlayerData.Tag + "   🚨",
		Description: "**" + matchData.Metadata.Map.Name + "**" + " -- " + gameOutcome + " -- **" + strconv.Itoa(roundsWon) + " : " + strconv.Itoa(roundsLost) + "**\n" + MMRData.Tier + " " + strconv.Itoa(MMRData.CurrentRR) + "RR (" + sign(MMRData.RRChange) + strconv.Itoa(MMRData.RRChange) + ")" + "\n\n" + getProgressBar(MMRData.CurrentRR),
		Fields:      embedFields,
		URL:         "https://tracker.gg/valorant/match/" + matchData.Metadata.MatchID,
		Color:       embedColor,
		Timestamp:   matchData.Metadata.StartedAt,
		Image: EmbedImage{
			URL:   maps[matchData.Metadata.Map.Name],
			Width: 1000,
		},
		Footer: EmbedFooter{
			Text: matchData.Metadata.MatchID,
		},
	}
	return embed
}

func getProgressBar(rating int) string {
	greenSquares := int((float64(rating) / 550) * 22)
	redSquares := 22 - greenSquares
	progressBar := []rune(strings.Repeat("🟩", greenSquares) + strings.Repeat("🟥", redSquares))
	progressBarRanks := append(progressBar[:3], []rune(getRankEmoji("Immortal2"))...)
	progressBarRanks = append(progressBarRanks, progressBar[4:7]...)
	progressBarRanks = append(progressBarRanks, []rune(getRankEmoji("Immortal3"))...)
	progressBarRanks = append(progressBarRanks, progressBar[8:21]...)
	progressBarRanks = append(progressBarRanks, []rune(getRankEmoji("Radiant"))...)
	return string(progressBarRanks)
}

func sign(x int) string {
	if x < 0 {
		return "-"
	}
	return "+"
}

func processMatchData(matchData MatchData, trackedPlayer Player) (int, int, int, string) {
	embedColor := 0
	gameState := 0
	roundsWon := 0
	roundsLost := 0

	for _, team := range matchData.Teams {

		playersTeam := false

		if team.TeamID == trackedPlayer.Team {
			playersTeam = true
			roundsWon = team.Rounds.Won
			roundsLost = team.Rounds.Lost
		}

		if team.Won {
			if playersTeam {
				gameState = 1
			} else {
				gameState = -1
			}
		}
	}

	var gameOutcome string
	switch gameState {
	case -1:
		gameOutcome = "**DEFEAT**"
		embedColor = LOSS_COLOR
	case 1:
		gameOutcome = "**VICTORY**"
		embedColor = WIN_COLOR
	case 0:
		gameOutcome = "**DRAW**"
		embedColor = DRAW_COLOR
	}
	return embedColor, roundsWon, roundsLost, gameOutcome
}

func executeWebhook(webhookURL string, matchData MatchData, trackedPlayerData TrackedPlayerData, MMRData MMRData) {
	var webhookData WebhookData = WebhookData{
		Embeds: []Embed{CreateEmbed(matchData, trackedPlayerData, MMRData)},
	}
	jsonData, err := json.Marshal(webhookData)
	if err != nil {
		log.Println("Error marshalling webhook data:", err)
	}

	res, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error making webhook request:", err)
	}
	defer res.Body.Close()
}

func getMatchData(client *http.Client, req *http.Request) MatchData {
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making HD API request:", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Println("Unexpected HD API response status:", res.Status)
	}

	matchDataResponse := handleRes(res)

	return matchDataResponse.MatchData[0]
}

func handleRes(res *http.Response) MatchDataResponse {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading HD API response body:", err)
	}

	var matchDataResponse MatchDataResponse

	err = json.Unmarshal(body, &matchDataResponse)

	if err != nil {
		log.Println("Error unmarshalling HD API response body:", err)
	}
	return matchDataResponse
}

func createUrl(playerData TrackedPlayerData) string {
	return "https://api.henrikdev.xyz/valorant/v4/matches/" + playerData.Region + "/" + playerData.Platform + "/" + playerData.Name + "/" + playerData.Tag + "?size=1&mode=competitive"
}

func main() {
	log.Println("START")
	log.Println("Loading environmental variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("Loading environmental variables...")
	WEBHOOK_URL := os.Getenv("WEBHOOK_URL")
	API_KEY := os.Getenv("API_KEY")

	var playerData TrackedPlayerData = TrackedPlayerData{
		Name:     os.Getenv("PLAYER_NAME"),
		Tag:      os.Getenv("PLAYER_TAG"),
		Platform: os.Getenv("PLAYER_PLATFORM"),
		Region:   os.Getenv("PLAYER_REGION"),
	}

	client := &http.Client{}

	apiUrl := createUrl(playerData)

	matchReq, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	matchReq.Header.Add("Authorization", API_KEY)

	log.Println("Server Started")
	var lastMatchID string

	for {
		log.Println("Checking match data...")
		matchData := getMatchData(client, matchReq)

		if matchData.Metadata.MatchID == lastMatchID {
			log.Println("No new match data found")
		} else {
			lastMatchID = matchData.Metadata.MatchID
			log.Println("New match found. Looking up MMR...")
			mmrApiUrl := "https://api.henrikdev.xyz/valorant/v3/mmr/" + playerData.Region + "/" + playerData.Platform + "/" + playerData.Name + "/" + playerData.Tag
			mmrReq, err := http.NewRequest("GET", mmrApiUrl, nil)
			if err != nil {
				log.Fatal(err)
			}
			mmrReq.Header.Add("Authorization", API_KEY)
			MMRData := getMMRData(client, mmrReq)
			executeWebhook(WEBHOOK_URL, matchData, playerData, MMRData)
		}

		time.Sleep(10 * time.Second)

	}
}

func getMMRData(client *http.Client, req *http.Request) MMRData {
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error making HD API request:", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading HD API response body:", err)
	}

	var mmrData MMRDataResponse
	json.Unmarshal(body, &mmrData)

	return MMRData{
		CurrentRR: mmrData.Data.Current.RR,
		RRChange:  mmrData.Data.Current.LastChange,
		Tier:      mmrData.Data.Current.Tier.Name,
	}
}
