package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type AppList struct {
	AppList ListOfGames `json:"applist"`
}

type ListOfGames struct {
	GameList []Game `json:"apps"`
}

type Game struct {
	// {"appid":1829051,"name":""}
	AppID int    `json:"appid"`
	Name  string `json:"name"`
}

const Environment = "prod" // "dev" or "prod", check GetEndpoint()

var SteamEndpoint string = GetEndpoint()

/*
	getRandomGame() : Give random game on Steam.
	getGameByName(name string) : Give the game matching the given name on Steam.
	getGameByGenre(genre SteamGenre) : Give a random game matching the given genre(s) on Steam.
*/

func GetRandomGame() (retour string, nb int) {
	appList := AppList{}
	resp, err := http.Get(SteamEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &appList)
		randomGame := ReturnRandomGame(appList)
		return fmt.Sprintf("https://store.steampowered.com/app/%d/",
			randomGame.AppID,
		), len(appList.AppList.GameList)
	} else {
		return "Nothing", 0
	}
}

func GetGameByName(name string) (retour string, nb int) {
	appList := AppList{}
	resp, err := http.Get(SteamEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &appList)

		for i, v := range appList.AppList.GameList {
			if v.Name == name {
				return fmt.Sprintf("https://store.steampowered.com/app/%d/",
					appList.AppList.GameList[i].AppID,
				), len(appList.AppList.GameList)
			}
		}
		return fmt.Sprintf("%s", name), len(appList.AppList.GameList)
	} else {
		return "Nothing", 0
	}
}

func GetGameByGenre(genre []string) (retour string, nb int) {
	appList := AppList{}
	genreEndpoint := SteamEndpoint + "?tag=" + strings.Join(genre, "&tag=")
	fmt.Print(genreEndpoint)
	resp, err := http.Get(genreEndpoint)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), &appList)
		randomGame := ReturnRandomGame(appList)
		return fmt.Sprintf("https://store.steampowered.com/app/%d/",
			randomGame.AppID,
		), len(appList.AppList.GameList)
	} else {
		return "There was a problem, here is the request tried : " + genreEndpoint, 0
	}
}

func ReturnRandomGame(appList AppList) (retour Game) {
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(len(appList.AppList.GameList))
	return appList.AppList.GameList[randomNumber]
}

func GetEndpoint() string {
	if Environment == "prod" {
		return "http://api.steampowered.com/ISteamApps/GetAppList/v0002/"
	}
	return "http://api.steampowered.com/ISteamApps/GetAppList/v2/" // smaller list, for test and dev
}
