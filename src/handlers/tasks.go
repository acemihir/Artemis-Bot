package handlers

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-co-op/gocron"
	"github.com/jerskisnow/Artemis-Bot/src/utils"
)

func FlushSuggestions() {
	var cursor uint64
	for {
		keys, cursor, ex := utils.Cache.Client.Scan(utils.Cache.Context, cursor, "prefix:s_", 0).Result()
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

			// TODO: Parse JSON

			_, ex = utils.Firebase.Firestore.Collection("submissions").Doc(k).Update(utils.Firebase.Context, []firestore.Update{
				{
					Path:  "upvotes",
					Value: upvote_array,
				},
				{
					Path:  "downvotes",
					Value: downvote_array,
				},
			})
			if ex != nil {
				utils.Cout("[ERROR] Could not update in Firestore: %v", utils.Red, ex)
				os.Exit(1) // Exit before more unsavable data will be collected
			}

			// Get values
			fmt.Println("key", k)
		}

		if cursor == 0 {
			break
		}
	}
}

func RegisterTasks() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(5).Hours().Do(FlushSuggestions)

	s.StartAsync()
}
