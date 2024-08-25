package state

import "gopkg.in/telebot.v3"

var (
	userCallbacks = make(map[int64]func(c telebot.Context) error)
)

func SetUserCallback(userID int64, callback func(c telebot.Context) error) {
	userCallbacks[userID] = callback
}

func DeleteUserCallback(userID int64) {
	delete(userCallbacks, userID)
}

func GetUserCallback(userID int64) func(c telebot.Context) error {
	callback, ok := userCallbacks[userID]
	if !ok {
		return nil
	}
	return callback
}
