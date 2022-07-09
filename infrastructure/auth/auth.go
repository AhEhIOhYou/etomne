package auth

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type AuthInterface interface {
	CreateAuth(uint64, *TokenDetails) error
	FetchAuth(string) (uint64, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type ClientData struct {
	client *redis.Client
}

func NewAuth(client *redis.Client) *ClientData {
	return &ClientData{
		client: client,
	}
}

type AccessDetails struct {
	TokenUuid string
	UserId    uint64
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func (tk *ClientData) CreateAuth(userId uint64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.TokenUuid, strconv.Itoa(int(userId)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tk.client.Set(td.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record instead")
	}
	return nil
}

func (tk *ClientData) FetchAuth(tokenUuid string) (uint64, error) {
	userId, err := tk.client.Get(tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userId, 10, 64)
	return userID, err
}

func (tk *ClientData) DeleteRefresh(refreshUuid string) error {
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted != 0 {
		return err
	}
	return nil
}

func (tk *ClientData) DeleteTokens(authD *AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%d", authD.TokenUuid, authD.UserId)
	deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	if err != nil {
		return err
	}
	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

var _ AuthInterface = &ClientData{}
