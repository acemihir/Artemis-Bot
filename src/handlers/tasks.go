package handlers

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/OnlyF0uR/Artemis-Bot/src/utils"
	"github.com/go-co-op/gocron"
)

func FlushSuggestions() {
	var cursor uint64
	for {
		keys, cursor, ex := utils.Cache.Client.Scan(utils.Cache.Context, cursor, "*s_*", 0).Result()
		if ex != nil {
			utils.Cout("[ERROR] Failed to scan Redis: %v", utils.Red, ex)
			os.Exit(1) // Exit before more unflashable data will be collected
		}

		for _, k := range keys {
			v, ex := utils.Cache.GetCache(k)
			if ex != nil {
				utils.Cout("[ERROR] Could not get from Redis: %v", utils.Red, ex)
				os.Exit(1) // Exit before more unflashable data will be collected
			}

			ex = utils.Cache.DelCache(k)
			if ex != nil {
				utils.Cout("[ERROR] Could not delete from Redis: %v", utils.Red, ex)
				continue
			}

			vote_data := utils.SuggestionVotes{}
			ex = json.Unmarshal([]byte(v), &vote_data)
			if ex != nil {
				utils.Cout("[ERROR] Could not parse JSON: %v", utils.Red, ex)
				os.Exit(1) // Exit before more unsavable data will be collected
			}

			// fmt.Println(vote_data)

			_, ex = utils.Firebase.Firestore.Collection("submissions").Doc(k).Update(utils.Firebase.Context, []firestore.Update{
				{
					Path:  "upvotes",
					Value: vote_data.Upvotes,
				},
				{
					Path:  "downvotes",
					Value: vote_data.Downvotes,
				},
			})
			if ex != nil {
				if strings.Contains(ex.Error(), "No document to update") {
					utils.Cout("[WARN] Firestore updated failed, entry not present.", utils.Yellow)
				} else {
					utils.Cout("[ERROR] Could not update in Firestore: %v", utils.Red, ex)
					os.Exit(1) // Exit before more unsavable data will be collected
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	utils.Cout("[INFO] Suggestions have been flushed.", utils.Cyan)
}

var scheduler *gocron.Scheduler

func RegisterTasks(isProd bool) {
	scheduler = gocron.NewScheduler(time.UTC)

	if isProd {
		scheduler.Every(5).Hours().Do(FlushSuggestions)
		utils.Cout("[INFO] Suggestions will be flushed every 5 hours.", utils.Cyan)
	} else {
		scheduler.Every(8).Minutes().Do(FlushSuggestions)
		utils.Cout("[INFO] Suggestions will be flushed every 8 minutes.", utils.Cyan)
	}

	scheduler.StartAsync()
}

func ShutdownTasks() {
	scheduler.RunAllWithDelay(time.Duration(30) * time.Second)
	scheduler.Clear() // Cleanup
}
