package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var data TrackingData

func main() {
	port := flag.String("p", "8080", "port")
	session := flag.String("s", "6949868c-b1f0-444d-a1d1-9e50b06dcd18", "session")
	pollInterval := flag.Int("i", 1, "interval sec")
	flag.Parse()
	log.Info("port:", *port)
	go poll(*session, *pollInterval)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		i := data.TrackPoints[len(data.TrackPoints)-1]
		c.Data(
			200,
			"text/csv",
			[]byte(
				fmt.Sprintf(
					"ts,lat,lon,sp,hr\n%v,%v,%v,%v,%v\n",
					i.DateTime,
					i.Position.Lat,
					i.Position.Lon,
					i.Speed,
					i.FitnessPointData.HeartRateBeatsPerMin,
				),
			),
		)
	})
	r.Run("localhost:" + *port)
}

func poll(session string, rate int) {
	for {
		tnow := time.Now().UnixNano() / 1e6
		u := fmt.Sprintf("https://livetrack.garmin.com/services/session/%v/trackpoints?requestTime=%v&from=%v", session, tnow, tnow-60*1e3)
		log.Info("URL:", u)
		r, err := http.Get(u)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Info(string(b))
		err = json.Unmarshal(b, &data)
		if err != nil {
			log.Error(err)
			continue
		}
		time.Sleep(time.Second * time.Duration(rate))
	}
}

type TrackingData struct {
	TrackPoints []struct {
		DateTime time.Time `json:"dateTime"`
		Position struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"position"`
		Altitude         float64 `json:"altitude"`
		Speed            float64 `json:"speed"`
		FitnessPointData struct {
			DeviceID             int64         `json:"deviceId"`
			ActivityType         string        `json:"activityType"`
			EventTypes           []interface{} `json:"eventTypes"`
			PointStatus          string        `json:"pointStatus"`
			DistanceMeters       float64       `json:"distanceMeters"`
			HeartRateBeatsPerMin int           `json:"heartRateBeatsPerMin"`
			SpeedMetersPerSec    float64       `json:"speedMetersPerSec"`
			CadenceCyclesPerMin  int           `json:"cadenceCyclesPerMin"`
			DurationSecs         int           `json:"durationSecs"`
			ActivityCreatedTime  time.Time     `json:"activityCreatedTime"`
			TotalDistanceMeters  float64       `json:"totalDistanceMeters"`
			TotalDurationSecs    int           `json:"totalDurationSecs"`
			Elevation            string        `json:"elevation"`
		} `json:"fitnessPointData"`
	} `json:"trackPoints"`
}
