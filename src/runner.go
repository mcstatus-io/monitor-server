package main

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func StartRunner() {
	for {
		time.Sleep(time.Until(time.Now().Add(config.CycleInterval).Truncate(config.CycleInterval)))

		servers, err := db.GetNextUniqueServers()

		if err != nil {
			log.Fatal(err)
		}

		for _, server := range servers {
			status, err := GetServerStatus(server)

			if err != nil {
				log.Println(err)

				continue
			}

			if status.Online {
				if err = db.UpdateUniqueServerByID(server.ID, bson.M{
					"$inc": bson.M{
						"onlineCount": 1,
						"totalCount":  1,
					},
					"$set": bson.M{
						"status":       "online",
						"lastPingedAt": time.Now().UTC(),
					},
				}); err != nil {
					log.Println(err)

					continue
				}
			} else {
				if err = db.UpdateUniqueServerByID(server.ID, bson.M{
					"$inc": bson.M{
						"totalCount": 1,
					},
					"$set": bson.M{
						"status":       "offline",
						"lastPingedAt": time.Now().UTC(),
					},
				}); err != nil {
					log.Println(err)

					continue
				}
			}
		}
	}
}
