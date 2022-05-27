package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"gophkeeper/internal/common/models"
	"gophkeeper/pkg/crypto"
	"io/ioutil"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type Storage struct {
	Data            []models.Note
	LastSyncDate    time.Time
	fileStoragePath string
	Check           bool
}

type producer struct {
	file   *os.File
	writer *bufio.Writer
}

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func newProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func newConsumer(fileName string) (*consumer, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file: f,
	}, nil
}

func New(fileStoragePath string) *Storage {
	data := []models.Note{}
	storage := Storage{
		Data:            data,
		fileStoragePath: fileStoragePath,
		Check:           true,
	}
	return &storage
}

func (s *Storage) AddItem(title string, body string, key []byte) error {
	encrypted, err := crypto.Encrypt(key, []byte(body))
	if err != nil {
		return err
	}
	note := models.Note{
		Title:     title,
		Body:      encrypted,
		CreatedAt: time.Now(),
	}
	s.Data = append(s.Data, note)
	return nil
}

func (s *Storage) GetNotes(all bool) []models.Note {
	ret := make([]models.Note, 0)
	for _, d := range s.Data {
		if all || !d.Deleted {
			ret = append(ret, d)
		}
	}
	return ret
}

func (s *Storage) GetNonSyncedData() []models.Note {
	data := make([]models.Note, 0)
	for _, d := range s.Data {
		if d.DeletedAt.After(s.LastSyncDate) || d.CreatedAt.After(s.LastSyncDate) || d.UpdatedAt.After(s.LastSyncDate) {
			data = append(data, d)
		}
	}
	return data
}

// Updatedata обновляет записи в коллекции, удаляет записи с 0 id, добавляет новые
func (s *Storage) UpdateData(newdata []models.Note, lastSync time.Time) error {
	filtered := make([]models.Note, 0)
	// удаление не синхронизированых и обновление синхронизированых
	for _, d := range s.Data {
		if d.ID == "" || d.Deleted {
			continue
		}
		touched := false
		for _, nd := range newdata {
			if nd.ID == d.ID {
				if !nd.Deleted {
					filtered = append(filtered, nd)
				}
				touched = true
			}
		}
		if !touched {
			filtered = append(filtered, d)
		}
	}
	// добавление новых
	for _, nd := range newdata {
		if nd.CreatedAt.After(s.LastSyncDate) && !nd.Deleted {
			filtered = append(filtered, nd)
		}
	}
	s.LastSyncDate = lastSync
	s.Data = filtered
	return nil
}

// Flush save storage to file
func (storage *Storage) Flush(hash []byte) error {
	var buff bytes.Buffer
	gobEncoder := gob.NewEncoder(&buff)
	gobEncoder.Encode(storage)
	encrypted, _ := crypto.Encrypt(hash, buff.Bytes())
	file, err := os.OpenFile(storage.fileStoragePath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Err(err).Msg("open file error")
		return err
	}
	writer := bufio.NewWriter(file)
	writer.Write(encrypted)
	file.Truncate(int64(writer.Buffered()))
	return writer.Flush()
}

// CheckFile check file presence
func (storage *Storage) CheckFile() bool {
	_, err := os.OpenFile(storage.fileStoragePath, os.O_RDONLY, 0777)
	if err == nil {
		return true
	}
	return false
}

// load local data storage from file
func (storage *Storage) Load(hash []byte) error {
	f, err := os.OpenFile(storage.fileStoragePath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	readFile, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	decrypted, err := crypto.Decrypt(hash, readFile)
	if err != nil {
		return err
	}
	strg := Storage{}
	var buff bytes.Buffer
	dec := gob.NewDecoder(&buff)
	buff.Write(decrypted)
	err = dec.Decode(&strg)
	if err != nil {
		return err
	}
	if strg.Check {
		storage.Data = strg.Data
		storage.LastSyncDate = strg.LastSyncDate
		return nil
	}
	return errors.New("data is corrupted!")
}

func (s *Storage) GetByIndex(index int) (models.Note, error) {
	if index > len(s.Data)-1 {
		return models.Note{}, errors.New("no such note")
	}
	return s.Data[index], nil
}

func (s *Storage) SetDeleted(index int) error {
	if index > len(s.Data)-1 {
		return errors.New("no such note")
	}
	if s.Data[index].ID == "" {
		if index > (len(s.Data) - 2) {
			s.Data = s.Data[:index]
		} else {
			s.Data = append(s.Data[:index], s.Data[index+1:]...)
		}
		return nil
	}
	s.Data[index].DeletedAt = time.Now()
	s.Data[index].Deleted = true
	return nil
}

func (s *Storage) GetLastSyncDate() time.Time {
	return s.LastSyncDate
}

func (s *Storage) GetDataLen() int {
	return len(s.Data)
}
