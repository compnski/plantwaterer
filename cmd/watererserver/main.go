package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

/*
{"now": 1485775 ,"sensors":[
 {"id": 0 ,"d":[ 1485771 , 37440 , 1485772 , 37376 , 1485773 , 37376 , 1485774 , 37440 , 1485775 , 37440 , 1485766 , 37312 , 1485767 , 37376 , 1485768 , 37248 , 1485769 , 37376 ]}
, {"id": 1 ,"d":[ 1485771 , 37376 , 1485772 , 37440 , 1485773 , 37376 , 1485774 , 37376 , 1485775 , 37376 , 1485766 , 37376 , 1485767 , 37376 , 1485768 , 37376 , 1485769 , 37376 ]}
, {"id": 2 ,"d":[ 1485771 , 37632 , 1485772 , 37632 , 1485773 , 37632 , 1485774 , 37632 , 1485775 , 37696 , 1485766 , 37632 , 1485767 , 37632 , 1485768 , 37632 , 1485769 , 37632 ]}
, {"id": 3 ,"d":[ 1485771 , 36608 , 1485772 , 36608 , 1485773 , 36608 , 1485774 , 36608 , 1485775 , 36608 , 1485766 , 36608 , 1485767 , 36608 , 1485768 , 36608 , 1485769 , 36608 ]}
],"sections":[
 {"id": 0 , "on": false , "next": 1555513 , "last": 1296313 , "onTime": 50 , "offTime": 259200 , "onAcc": 306 , "offAcc": 1296007 }
]}
*/

//, "onTime": 50 , "offTime": 259200
//, "onAcc": 306 , "offAcc": 1296007 }
type SensorSection struct {
	Id   int
	Data []int `json:"d"`
}

func (s SensorSection) DataAsMap() map[int]int {
	parsedData := map[int]int{}
	for i := 0; i < len(s.Data)-1; i += 2 {
		at := s.Data[i]
		value := s.Data[i+1]
		parsedData[at] = value
	}
	return parsedData
}

type t struct{}

func printSensorsAsCsv(w *csv.Writer, timeOffset uint, sensors []SensorSection, startTs int, writeHeader bool) (lastTs int) {
	// get all timestamps
	tsSet := map[int]t{}
	tsList := sort.IntSlice{}
	for _, s := range sensors {
		for ts := range s.DataAsMap() {
			if _, exists := tsSet[ts]; !exists && ts > startTs {
				tsSet[ts] = t{}
				tsList = append(tsList, int(ts))
			}
		}
	}
	sort.Sort(tsList)
	rec := make([]string, len(sensors)+1)
	if writeHeader {
		rec[0] = "ts"
		for idx, s := range sensors {
			rec[idx+1] = fmt.Sprintf("sensor_%d", s.Id)
		}
		w.Write(rec)
	}
	for _, ts := range tsList {
		rec[0] = strconv.Itoa(ts + int(timeOffset))
		for idx, s := range sensors {
			rec[idx+1] = strconv.Itoa(int(s.DataAsMap()[int(ts)]))
		}
		w.Write(rec)
	}
	if len(tsList) == 0 {
		return startTs
	}
	return tsList[len(tsList)-1]
}

func printSectionsAsCsv(w *csv.Writer, timeOffset uint, sections []WaterSection, ts int, writeHeader bool) (lastTs int) {
	rec := make([]string, len(sections)+1)
	if writeHeader {
		rec[0] = "ts"
		for idx, s := range sections {
			rec[idx+1] = fmt.Sprintf("section_%d", s.Id)
		}
		w.Write(rec)
	}
	rec[0] = strconv.Itoa(ts + int(timeOffset))
	for idx, s := range sections {
		rec[idx+1] = strconv.Itoa(s.OnAcc)
	}
	w.Write(rec)
	return ts
}

type WaterSection struct {
	Id      int
	On      bool
	Next    int
	Last    int
	OnTime  int
	OffTime int
	OnAcc   int
	OffAcc  int
}

type StatusMessage struct {
	Now      int
	Sensors  []SensorSection
	Sections []WaterSection
}

func (m StatusMessage) Valid() bool {
	return m.Now > 0
}

type CommandResponseMessage struct {
	Command int `json:"cmd"`
	Id      int
	Value   int
	Status  int
}

func (m CommandResponseMessage) Valid() bool {
	return m.Command > 0
}

type CommandErrorMessage struct {
	Baddata string
}

type RecordUnion struct {
	StatusMessage
	CommandResponseMessage
	CommandErrorMessage
}

func main() {
	var record RecordUnion
	var timeOffset = flag.Uint("offset", 1602157508, "Epoch offset to add to Ts to match real epoch time")
	flag.Parse()
	dec := json.NewDecoder(os.Stdin)
	sensorFile, err := os.OpenFile("sensor.csv", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer sensorFile.Close()
	sensorW := csv.NewWriter(sensorFile)

	waterFile, err := os.OpenFile("water.csv", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer waterFile.Close()
	waterW := csv.NewWriter(waterFile)
	defer waterW.Flush()
	defer sensorW.Flush()

	var lastTs int
	for {
		err := dec.Decode(&record)
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal(dec.InputOffset(), err)
		}
		if record.StatusMessage.Valid() {
			// Write headers on the first record
			printSectionsAsCsv(waterW, *timeOffset, record.Sections, record.StatusMessage.Now, lastTs == 0)
			lastTs = printSensorsAsCsv(sensorW, *timeOffset, record.Sensors, lastTs, lastTs == 0)
		} else {
			fmt.Printf("%+v\n\n", record)
		}
	}

}
