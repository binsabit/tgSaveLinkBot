package files

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"

	"github.com/binsabit/tgSaveLinkBot/storage"
)

type Storage struct {
	basepath string
}

func New(basepath string) *Storage {
	return &Storage{
		basepath: basepath,
	}
}

func (s *Storage) Insert(link storage.Link) error {
	_, err := os.Stat(path.Join(s.basepath, strconv.Itoa(link.ChatId)))
	if os.IsNotExist(err) {
		err = os.Mkdir(path.Join(s.basepath, strconv.Itoa(link.ChatId)), os.ModePerm)
		if err != nil {
			return err
		}
	}

	hash := link.Hash()
	filepath := path.Join(s.basepath, strconv.Itoa(link.ChatId), hash)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	_, _ = file.WriteString(link.Content)
	defer file.Close()

	return nil
}

func (s *Storage) Delete(link storage.Link) error {
	hash := link.Hash()
	filepath := path.Join(s.basepath, strconv.Itoa(link.ChatId), hash)
	if !s.IsExists(link) {
		return errors.New("do not contain such file")
	}

	err := os.Remove(filepath) // remove a single file
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetRandom(chatId int) (storage.Link, error) {

	files, err := ioutil.ReadDir(path.Join(s.basepath, strconv.Itoa(chatId)))
	if err != nil {
		return storage.Link{}, err
	}
	if len(files) == 0 {
		return storage.Link{}, storage.ErrNoLinkFound
	}

	random := rand.Intn(len(files))
	content, err := os.ReadFile(path.Join(s.basepath, strconv.Itoa(chatId), files[random].Name()))
	if err != nil {
		return storage.Link{}, err
	}

	return storage.Link{
		ChatId:  chatId,
		Content: string(content),
	}, nil

}

func (s *Storage) IsExists(link storage.Link) bool {
	hash := link.Hash()
	filepath := path.Join(s.basepath, hash)
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
