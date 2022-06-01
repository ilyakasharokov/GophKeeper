// storage keeps data in memory and save it to disk
package storage

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/models"
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

func New(fileStoragePath string) *Storage {
	data := []models.Note{}
	storage := Storage{
		Data:            data,
		fileStoragePath: fileStoragePath,
		Check:           true,
	}
	return &storage
}

// AddNote to local storage
func (s *Storage) AddNote(title string, body string, key []byte) error {
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

// GetNotes gets all notes or only non-deleted notes
func (s *Storage) GetNotes(all bool) []models.Note {
	ret := make([]models.Note, 0)
	for _, d := range s.Data {
		if all || !d.Deleted {
			ret = append(ret, d)
		}
	}
	return ret
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

// Load local data storage from file
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

// GetNote get note by index
func (s *Storage) GetNote(index int) (models.Note, error) {
	if index > len(s.Data)-1 {
		return models.Note{}, errors.New("no such note")
	}
	return s.Data[index], nil
}

// DeleteNode delete local node or set node deleted if it was synchronized
func (s *Storage) DeleteNote(index int) error {
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

// GetLastSyncDate get last synchronization data
func (s *Storage) GetLastSyncDate() time.Time {
	return s.LastSyncDate
}

// SetLastSyncDate set last synchronization data
func (s *Storage) SetLastSyncDate(date time.Time) {
	s.LastSyncDate = date
}

// GetNotesCount get count of notes in the storage
func (s *Storage) GetNotesCount() int {
	return len(s.Data)
}

// SetNotes reset storage with new data
func (s *Storage) SetNotes(notes []models.Note) {
	s.Data = notes
}
