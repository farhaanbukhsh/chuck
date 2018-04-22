package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const chuckAPI = "http://api.chucknorris.io/jokes/random"

var home = os.Getenv("HOME")
var dbLocation = home + "/.jokesdb"

// Chucknorris is the struct used to  unmarshal the JSON response from the URL
type Chucknorris struct {
	Category []string `json:"category"`
	IconURL  string   `json:"icon_url"`
	ID       string   `json:"id"`
	URL      string   `json:"url"`
	Value    string   `json:"value"`
}

/*  getJokes takes the API url as the parameter and fetch jokes from it,
 *   here we have assumed it to be of ChuckNorris type
 */

//TODO: Remove the dependency on the hardcoded struct
func getJokes(URL string) (string, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", fmt.Errorf("No request formed %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("No response: %v", err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Read error")
	}

	var joke Chucknorris
	if err = json.Unmarshal(respData, &joke); err != nil {
		return "", fmt.Errorf("Error in unmarsheling, %v", err)
	}

	return joke.Value, nil
}

// Check if the Database exist if it does, then flushes it off
func deleteExistingDatabase(name string) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		if err := os.Remove(name); err != nil {
			return fmt.Errorf("Can't Delete file: %v", err)
		}
	}
	return nil
}

/*
 * cacheUpJokes queries for the joke and stores it in the Database, the parameter it is,
 * numberOfJokes, i.e the number of jokes that has to be stored, every joke has a joke id(jid)
 * and the joke itself, this is stored in the db.
 */
func cacheUpJokes(numberOfJokes int) error {

	if err := deleteExistingDatabase(dbLocation); err != nil {
		return fmt.Errorf("Cannot delete old database, %v", err)
	}

	db, err := sql.Open("sqlite3", dbLocation)
	if err != nil {
		return fmt.Errorf("Couldn't make DB, %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE jokes (jid INTEGER PRIMARY KEY, joke VARCHAR(64) NULL)")
	if err != nil {
		return fmt.Errorf("Table can't be created, %v", err)
	}

	stmt, err := db.Prepare("INSERT INTO jokes(jid, joke) values(?,?)")
	if err != nil {
		return fmt.Errorf("Insert statement couldn't be formed: %v", err)
	}

	for i := 1; i <= numberOfJokes; i++ {
		joke, err := getJokes(chuckAPI)
		if err != nil {
			return fmt.Errorf("Jokes couldn't be fetched, %v", err)
		}

		_, err = stmt.Exec(i, joke)
		if err != nil {
			return fmt.Errorf("Could'nt make record entry: %v", err)
		}
	}
	return nil
}

func fetchJoke() (string, error) {
	db, err := sql.Open("sqlite3", dbLocation)
	defer db.Close()
	var count int
	var joke string

	if err != nil {
		return "", fmt.Errorf("Couldn't make DB, %v", err)
	}

	totalRow, err := db.Query("SELECT count(*) FROM jokes")
	if err != nil {
		return "", fmt.Errorf("Can't get number of rows, is the db created run chuck --index=5")
	}

	for totalRow.Next() {
		err = totalRow.Scan(&count)
		if err != nil {
			return "", fmt.Errorf("Rows can't be read")
		}
	}

	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(count-1) + 1

	stm, err := db.Prepare("SELECT joke FROM jokes where jid = ?")
	if err != nil {
		return "", fmt.Errorf("Can't Prepaer statement")
	}

	res, err := stm.Query(randNum)
	if err != nil {
		return "", fmt.Errorf("Not able to fetch jokes")
	}

	for res.Next() {
		if err = res.Scan(&joke); err != nil {
			return "", fmt.Errorf("No jokes found")
		}
	}

	return joke, nil
}

func main() {
	index := flag.Int("index", -1, "To cache up Chuck facts, the parameter decide the number of jokes to cache")
	flag.Parse()
	if *index > 0 {
		if err := cacheUpJokes(*index); err != nil {
			log.Fatalf("Database issue: %v", err)
			os.Exit(2)
		}
	} else {
		joke, err := fetchJoke()
		if err != nil {
			fmt.Printf("No jokes found did you cache it by chuck --index=5 %v", err)
			os.Exit(2)
		}
		fmt.Println(joke)
	}
}
