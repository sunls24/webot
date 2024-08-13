package cache

import (
	"sync"
	"webot/internal/constants"
	"webot/internal/database"
	"webot/internal/entity"
	"webot/internal/types"
)

var cache = make(map[string]map[types.SettingsKey]string)
var lock sync.RWMutex

func InitSettings(db *database.DB) {
	lock.Lock()
	defer lock.Unlock()
	var list []entity.Settings
	if err := db.Find(&list).Error; err != nil {
		return
	}
	for _, v := range list {
		s, ok := cache[v.AvatarID]
		if !ok {
			s = make(map[types.SettingsKey]string)
			cache[v.AvatarID] = s
		}
		s[types.SettingsKey(v.Key)] = v.Value
	}
}

func GetSettings(id string, key types.SettingsKey) string {
	lock.RLock()
	defer lock.RUnlock()
	s, ok := cache[id]
	if !ok {
		return constants.ON
	}
	return s[key]
}

func SetSettings(id string, key types.SettingsKey, value string, db *database.DB) {
	lock.Lock()
	defer lock.Unlock()
	s, ok := cache[id]
	if !ok {
		s = make(map[types.SettingsKey]string)
		cache[id] = s
	}
	s[key] = value

	go func() {
		var settings = entity.Settings{
			AvatarID: id,
			Key:      string(key),
		}
		err := db.Where(&settings).FirstOrCreate(&settings).Error
		if err != nil {
			return
		}
		settings.Value = value
		db.Save(&settings)
	}()
}
