package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"

	log "github.com/astaxie/beego/logs"
	_ "github.com/bmizerany/pq"
)

type session struct {
	token     string
	loginTime time.Time
}

var (
	sessionsLock   sync.RWMutex
	sessions       map[int64]*session = make(map[int64]*session)
	sessionTimeout time.Duration
)

func CheckSession(userID int64, token string) error {
	sessionsLock.RLock()
	defer sessionsLock.RUnlock()
	s, ok := sessions[userID]
	if !ok {
		log.Error("no login user:", userID)
		return errors.New("no login user")
	}
	if time.Now().Sub(s.loginTime) > sessionTimeout {
		log.Error(userID, "login session timeout")
		return errors.New("login session timeout")
	}
	if s.token != token {
		log.Error(userID, "token error:", token)
		return errors.New("token error")
	}

	return nil
}

func AddSession(userID int64) string {

	randBytes := make([]byte, 32)
	rand.Read(randBytes)
	token := base64.RawStdEncoding.EncodeToString(randBytes)

	sessionsLock.RLock()
	s, ok := sessions[userID]
	sessionsLock.RUnlock()
	if !ok {
		s = &session{
			loginTime: time.Now(),
			token:     token,
		}
		sessionsLock.Lock()
		sessions[userID] = s
		sessionsLock.Unlock()
		return token
	}
	sessionsLock.Lock()
	s.loginTime = time.Now()
	s.token = token
	sessionsLock.Unlock()
	return token
}

func DelSession(userID int64) {
	sessionsLock.Lock()
	delete(sessions, userID)
	sessionsLock.Unlock()
}
