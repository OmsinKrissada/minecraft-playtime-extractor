package main

import (
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Tnze/go-mc/nbt"
)

const TICKS_IN_AN_HOUR = 20 * 60 * 60

type UserCache struct {
	Name string
	Uuid string
}

func getUsernameMap() map[string]string {
	path := filepath.Join(config.ServerDir, "usercache.json")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var decoded []UserCache
	json.Unmarshal(data, &decoded)

	uuid_username_map := make(map[string]string)
	for _, uc := range decoded {
		uuid_username_map[uc.Uuid] = uc.Name
	}

	return uuid_username_map
}

func getSinglePlaytime(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	type DecodedType struct {
		Stats struct {
			Custom struct {
				Playtime int `json:"minecraft:play_time"`
			} `json:"minecraft:custom"`
		}
	}
	var decoded DecodedType
	json.Unmarshal(data, &decoded)
	pt := decoded.Stats.Custom.Playtime

	return pt, nil
}

func getSingleLastSeen(path string) (time.Time, error) {
	filestat, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return filestat.ModTime(), nil
}

func getAllPlaytime() (map[string]int, error) {
	path := filepath.Join(config.WorldDir, "stats")
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	user_playtime_tick := make(map[string]int)
	for _, f := range files {
		playerStatPath := filepath.Join(path, f.Name())
		pt, err := getSinglePlaytime(playerStatPath)
		if err != nil {
			return nil, err
		}
		username, _, _ := strings.Cut(f.Name(), ".")
		user_playtime_tick[username] = pt
	}

	return user_playtime_tick, nil
}

func getAllLastSeen() (map[string]time.Time, error) {
	path := filepath.Join(config.WorldDir, "stats")
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	player_lastseen := make(map[string]time.Time)
	for _, f := range files {
		playerStatPath := filepath.Join(path, f.Name())
		ls, err := getSingleLastSeen(playerStatPath)
		if err != nil {
			return nil, err
		}
		username, _, _ := strings.Cut(f.Name(), ".")
		player_lastseen[username] = ls
	}

	return player_lastseen, nil
}

type WorldNBT struct {
	Data struct {
		Time int64
	}
}

func getWorldRunTime() int64 {
	path := filepath.Join(config.WorldDir, "level.dat")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	var decoded WorldNBT
	nbtDecoder := nbt.NewDecoder(gzipReader)
	nbtDecoder.Decode(&decoded)
	return decoded.Data.Time
}
