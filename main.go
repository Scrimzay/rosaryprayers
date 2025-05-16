package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"math/rand"

	"github.com/gin-gonic/gin"
)
	
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("**/*.html")
	r.Static("/static", "./static")

	r.GET("/", indexHandler)
	r.GET("/:day", todayRosaryHandler)
	r.GET("/weekday/:day", RosaryDayHandler)
	r.GET("/void", func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "rosaryinfo.html", gin.H{
			"Warning": "The rosary is lost in the void.",
		})
	})

	err := r.Run(":4000")
	if err != nil {
		log.Fatal(err)
	}
}

// for displaying the audio files at each endpoint
var mp3Links = map[string]string{
	"Sorrowful": "/static/rosary_sorrowful.mp3",
	"Glorious": "/static/rosary_glorious.mp3",
	"Joyful": "/static/rosary_joyful.mp3",
	"Luminous": "/static/rosary_luminous.mp3",
}
	
func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func RosaryDayHandler(c *gin.Context) {
    // Get day from POST form or URL param
    day := c.PostForm("day")
    if day == "" {
        day = c.Param("day")
    }
    if day == "" {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{
            "Error": "No day specified. Choose a day to proceed.",
        })
        return
    }

    // Fetch API data
    info, err := DailyRosaryAPILogic(day)
    if err != nil {
        fmt.Printf("Could not call RosaryDay func for %s: %v\n", day, err)
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Error": "The library is haunted and cannot respond. Try again later.",
        })
        return
    }

    // Set default MP3 link or use map
    prayerLink := ""
    if link, exists := mp3Links[info.GroupBy]; exists {
        prayerLink = link
    } else {
        prayerLink = "/static/default_rosary.mp3" // Fallback or empty string
    }

	// Mystery maps
	JoyfulMysteries := map[string]string{
		"Announce1": "The Annunciation of the Angel to Mary",
		"Announce2": "The Visitation of Mary to Saint Elizabeth",
		"Announce3": "The Nativity of Jesus in Bethlehem",
		"Announce4": "The Presentation of Jesus in the Temple",
		"Announce5": "The Finding of Jesus in the Temple",
	}

	SorrowfulMysteries := map[string]string{
		"Announce1": "The Agony of Jesus in the Garden",
		"Announce2": "The Scourging of Jesus at the Pillar",
		"Announce3": "The Crowning with Thorns",
		"Announce4": "The Carrying of the Cross",
		"Announce5": "The Crucifixion and Death of Jesus",
	}

	GloriousMysteries := map[string]string{
		"Announce1": "The Resurrection of Jesus Christ",
		"Announce2": "The Ascension of Jesus to Heaven",
		"Announce3": "The Descent of the Holy Spirit",
		"Announce4": "The Assumption of the Blessed Virgin Mary into Heaven",
		"Announce5": "The Coronation of the Blessed Virgin Mary, Queen of Heaven and Earth",
	}

	LuminousMysteries := map[string]string{
		"Announce1": "The Baptism of Jesus in the Jordan",
		"Announce2": "The Wedding Feast at Cana",
		"Announce3": "The Proclamation of the Kingdom of God",
		"Announce4": "The Transfiguration",
		"Announce5": "The Institution of the Eucharist",
	}

	// Map GroupBy to mystery map
	groupByToMysteries := map[string]map[string]string{
		"Joyful":    JoyfulMysteries,
		"Sorrowful": SorrowfulMysteries,
		"Glorious":  GloriousMysteries,
		"Luminous":  LuminousMysteries,
	}

	// Select mystery map based on GroupBy
    mysteries := groupByToMysteries[info.GroupBy]
    if mysteries == nil {
        mysteries = JoyfulMysteries // Fallback to Joyful if GroupBy is unknown
    }

    // 5% chance of an anomaly
    rand.Seed(time.Now().UnixNano())
    if rand.Float64() < 0.05 {
        anomaly := rand.Intn(4)
        switch anomaly {
        case 0: // Corrupted MP3 Link
            prayerLink = "forbidden://void/█▒█-corrupted.mp3"
            c.HTML(http.StatusOK, "rosarybyday.html", gin.H{
                "ID":            info.ID,
                "RosaryDay":     info.Date,
                "Mystery":       info.GroupBy,
                "InTheName":     info.InTheName,
                "IBelieve":      info.IBelieve,
                "OurFather":     info.OurFather,
                "HailMary":      info.HailMary,
                "GloryBe":       info.GloryBe,
                "OhMyJesus":     info.OhMyJesus,
                "Announce":      JoyfulMysteries["Announce1"],
                "Prayer":        prayerLink,
                "Warning":       "The library's spirit has tampered with this audio...",
            })
            return
        case 1: // Haunted Message
            time.Sleep(3 * time.Second)
            c.HTML(http.StatusOK, "rosarybyday.html", gin.H{
                "ID":            info.ID,
                "RosaryDay":     info.Date,
                "Mystery":       info.GroupBy,
                "InTheName":     info.InTheName,
                "IBelieve":      info.IBelieve,
                "OurFather":     info.OurFather,
                "HailMary":      info.HailMary,
                "GloryBe":       info.GloryBe,
                "OhMyJesus":     info.OhMyJesus,
                "Announce":      info.Announce,
                "Prayer":        prayerLink,
                "Curse":         "Do not listen to this alone...",
            })
            return
        case 2: // Partial Render
            c.HTML(http.StatusOK, "rosarybyday.html", gin.H{
                "ID":        info.ID,
                "RosaryDay": info.Date,
                "Warning":   "The library has sealed this knowledge.",
            })
            return
        case 3: // Distorted Mystery
            info.GroupBy = "Forbidden"
            prayerLink = ""
            c.HTML(http.StatusOK, "rosarybyday.html", gin.H{
                "ID":            info.ID,
                "RosaryDay":     info.Date,
                "Mystery":       info.GroupBy,
                "InTheName":     info.InTheName,
                "IBelieve":      info.IBelieve,
                "OurFather":     info.OurFather,
                "HailMary":      info.HailMary,
                "GloryBe":       info.GloryBe,
                "OhMyJesus":     info.OhMyJesus,
                "Announce":      info.Announce,
                "Prayer":        prayerLink,
                "Warning":       "This mystery is not meant to be known.",
            })
            return
        }
    }

    // Normal response
    c.HTML(http.StatusOK, "rosarybyday.html", gin.H{
        "ID":            info.ID,
        "RosaryDay":     info.Date,
        "Mystery":       info.GroupBy,
        "InTheName":     info.InTheName,
        "IBelieve":      info.IBelieve,
        "OurFather":     info.OurFather,
        "HailMary":      info.HailMary,
        "GloryBe":       info.GloryBe,
        "OhMyJesus":     info.OhMyJesus,
        "Announce1":      mysteries["Announce1"],
        "Announce2":      mysteries["Announce2"],
        "Announce3":      mysteries["Announce3"],
        "Announce4":      mysteries["Announce4"],
        "Announce5":      mysteries["Announce5"],
		"HolyQueen": info.HolyQueen,
		"FinalPrayer": info.FinalPrayer,
        "Prayer":        prayerLink,
    })
}

func todayRosaryHandler(c *gin.Context) {
	day := c.PostForm("day")
    if day == "" {
        day = c.Param("day")
    }
    if day == "" {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{
            "Error": "No day specified. Choose a day to proceed.",
        })
        return
    }

	info, err := RosaryTodayAPILogic(day)
	if err != nil {
        fmt.Printf("Could not call RosaryToday func: %v\n", err)
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Error": "The library is haunted and cannot respond. Try again later.",
        })
        return
    }
	
	// Set default MP3 link or use map
    prayerLink := ""
    if link, exists := mp3Links[info.Mystery]; exists {
        prayerLink = link
    } else {
        prayerLink = "/static/default_rosary.mp3" // Fallback or empty string
    }

    // 5% chance of an anomaly
    rand.Seed(time.Now().UnixNano())
    if rand.Float64() < 0.05 {
        anomaly := rand.Intn(3)
        switch anomaly {
        case 0: // Corrupted MP3 Link
            prayerLink = "forbidden://void/█▒█-corrupted.mp3"
            c.HTML(http.StatusOK, "rosaryinfo.html", gin.H{
                "ID":          info.ID,
                "RosaryDate":  info.RosaryDate,
                "RosaryDay":   info.RosaryDay,
                "Mystery":     info.Mystery,
                "Season":      info.Season,
                "CurrentDate": info.CurrentDate,
                "Prayer":      prayerLink,
                "Warning":     "A demon is tampering with this prayer...",
            })
            return
        case 1: // Haunted Message
            time.Sleep(3 * time.Second)
            c.HTML(http.StatusOK, "rosaryinfo.html", gin.H{
                "ID":          info.ID,
                "RosaryDate":  info.RosaryDate,
                "RosaryDay":   info.RosaryDay,
                "Mystery":     info.Mystery,
                "Season":      info.Season,
                "CurrentDate": info.CurrentDate,
                "Prayer":      prayerLink,
                "Curse":       "Do not listen to this alone...",
            })
            return
        case 2: // Partial Render
            c.HTML(http.StatusOK, "rosaryinfo.html", gin.H{
                "ID":         info.ID,
                "RosaryDate": info.RosaryDate,
                // Stop rendering other fields
                "Warning": "A demon has sealed this knowledge.",
            })
            return
        }
    }

    // Normal response
    c.HTML(http.StatusOK, "rosaryinfo.html", gin.H{
        "ID":          info.ID,
        "RosaryDate":  info.RosaryDate,
        "RosaryDay":   info.RosaryDay,
        "Mystery":     info.Mystery,
        "Season":      info.Season,
        "CurrentDate": info.CurrentDate,
        "Prayer":      prayerLink,
    })
}

type YesterdayTodayTomorrowAPIResponse struct {
	ID string `json:"id"`
	RosaryDate string `json:"rosary_date"`
	RosaryDay string `json:"rosary_day"`
	Mystery string `json:"mystery"`
	Season string `json:"season"`
	CurrentDate string `json:"currentDate"`
}

// handles today, yesterday and tomorrow
func RosaryTodayAPILogic(day string) (*YesterdayTodayTomorrowAPIResponse, error) {
	baseURL := fmt.Sprintf("https://the-rosary-api.vercel.app/v1/%s", day)

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, fmt.Errorf("Could not send get to Rosary API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read resp from Rosary API: %v", err)
	}

	var values []YesterdayTodayTomorrowAPIResponse
	err = json.Unmarshal(body, &values)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal json data from Rosary API: %v", err)
	}

	return &values[0], nil
}

type DailyRosaryAPIResponse struct {
	ID            string `json:"id"`
    GroupBy       string `json:"group_by"` // Maps to Mystery
    Date          string `json:"date"`     // Maps to RosaryDay
    InTheName     string `json:"in_the_name_1"`
    IBelieve      string `json:"i_believe"`
    OurFather     string `json:"our _father_1"`
    HailMary      string `json:"hail _mary_1"`
    HailMaryCount int    // Not in JSON, set manually
    GloryBe       string `json:"glory_be_1"`
    OhMyJesus     string `json:"oh_my_jesus_1"`
    Announce      string `json:"announce_1"`
	HolyQueen string `json:"holy_queen"`
	FinalPrayer string `json:"final_prayer"`
}

// handles the general weekdays
func DailyRosaryAPILogic(day string) (*DailyRosaryAPIResponse, error) {
	baseURL := fmt.Sprintf("https://the-rosary-api.vercel.app/v1/%s", day)

	resp, err := http.Get(baseURL)
    if err != nil {
        return nil, fmt.Errorf("could not send GET to Rosary API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Rosary API returned non-200 status: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("could not read response from Rosary API: %v", err)
    }

    var value []DailyRosaryAPIResponse
    err = json.Unmarshal(body, &value)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshal JSON from Rosary API: %v, body: %s", err, string(body))
    }

	return &value[0], nil
}