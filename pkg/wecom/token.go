package wecom

import (
	"context"
	"sync"
	"time"
)

type TokenManager interface {
	Get(context.Context) (TokenInfo, error)
	Refresh(context.Context) (TokenInfo, error)
}

type TokenInfo struct {
	Token     string
	ExpiresIn time.Duration
}

type token struct {
	mutex *sync.RWMutex
	TokenInfo
	lastRefresh  time.Time
	getTokenFunc func() (TokenInfo, error)
}

func (t *token) setTokenFunc(f func() (TokenInfo, error)) {
	t.getTokenFunc = f
}

func (t *token) Get(ctx context.Context) (TokenInfo, error) {
	t.mutex.RLock()
	if t.Token == "" || time.Now().Sub(t.lastRefresh) >= t.ExpiresIn {
		t.mutex.RUnlock()
		_ = t.refreshToken()
		t.mutex.RLock()
	}
	defer t.mutex.RUnlock()
	return t.TokenInfo, nil
}

func (t *token) Refresh(ctx context.Context) (TokenInfo, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	err := t.refreshToken()
	if err != nil {
		return TokenInfo{}, err
	}
	return t.TokenInfo, nil
}

func (t *token) refreshToken() error {
	info, err := t.getTokenFunc()
	if err != nil {
		return err
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.Token = info.Token
	t.ExpiresIn = info.ExpiresIn * time.Second
	t.lastRefresh = time.Now()
	return nil
}

func (app *WeComApp) getToken() (TokenInfo, error) {
	t, err := app.getAccessToken(context.Background())
	if err != nil {
		return TokenInfo{}, err
	}
	return TokenInfo{
		Token:     t.AccessToken,
		ExpiresIn: time.Duration(t.ExpiresIn),
	}, nil
}
