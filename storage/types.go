package storage

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

type Storage interface {
	Insert(Link) error
	Delete(Link) error
	GetRandom(int) (Link, error)
	IsExists(Link) bool
}

var (
	ErrNoLinkFound = errors.New("no links found")
)

type Link struct {
	ChatId  int    `json:"chat_id"`
	Content string `json:"content"`
}

func New(content string, chatId int) Link {
	return Link{
		ChatId:  chatId,
		Content: content,
	}
}

func (l Link) Hash() string {
	h := sha256.New()
	h.Write([]byte(l.Content + strconv.Itoa(l.ChatId)))
	return l.ToBase64(h.Sum(nil))
}

func (l Link) ToBase64(data []byte) string {
	return strings.Replace(base64.StdEncoding.EncodeToString(data), "/", "", -1)
}
