// cli - command line interface of the app
package cli

import (
	"context"
	"errors"
	"fmt"

	"gophkeeper/internal/common/models"
	"gophkeeper/pkg/crypto"
	"regexp"
	"strconv"
	"time"
)

const (
	START = iota
	LOGGED_IN
)

type CLI struct {
	state   int
	client  GRPCClientModel
	storage StorageModel
	user    *models.User
}

type GRPCClientModel interface {
	Login(login string, pwd string) (string, error)
	Registration(login string, pwd string) (string, string, error)
	SyncData(notes []models.Note, lastSync time.Time) ([]models.Note, time.Time, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
	Close()
}

type StorageModel interface {
	AddItem(title string, body string, key []byte) error
	GetNotes(all bool) []models.Note
	GetNonSyncedData() []models.Note
	UpdateData(newData []models.Note, lastSync time.Time) error
	Flush(hash []byte) error
	CheckFile() bool
	Load(hash []byte) error
	GetByIndex(index int) (models.Note, error)
	SetDeleted(index int) error
	GetLastSyncDate() time.Time
	GetDataLen() int
}

var deleteCmdRegExp = regexp.MustCompile(`^d(\d+)$`)
var viewCmdRegExp = regexp.MustCompile(`^v(\d+)$`)

// registration register and return user if success
func (cli *CLI) registration() (models.User, error) {
	login, pwd := readAuth()
	status, token, err := cli.client.Registration(login, pwd)
	if err != nil {
		return models.User{}, err
	}
	if token == "" {
		return models.User{}, errors.New(status)
	}
	hash := crypto.Hash(pwd)
	return models.User{Login: login, PasswordHash: hash, Token: token}, nil
}

// auth authenticate and return user if success
func (cli *CLI) auth() (models.User, error) {
	login, pwd := readAuth()
	hash := crypto.Hash(pwd)
	var err error
	var token string
	if cli.client != nil {
		token, err = cli.client.Login(login, pwd)
		if err != nil {
			return models.User{}, errors.New("can't auth remotely")
		}
	}
	if cli.storage.CheckFile() {
		err = cli.storage.Load(hash)
		if err != nil {
			return models.User{}, err
		}
	} else {
		newData, newLastSync, err := cli.client.SyncData(nil, time.Time{})
		if err != nil {
			fmt.Println("sync error (" + err.Error() + ")")
		}
		err = cli.storage.UpdateData(newData, newLastSync)
		if err != nil {
			fmt.Println("merge data error (" + err.Error() + ")")
		}
		err = cli.storage.Flush(hash)
		if err != nil {
			fmt.Println("save notes error (" + err.Error() + ")")
		}
	}
	return models.User{Login: login, PasswordHash: hash, Token: token}, nil

}

// readAuth read login and password from command line
func readAuth() (login string, pwd string) {
	login, pwd = "", ""
	fmt.Print("Enter your name:")
	login = readString()
	fmt.Print("Enter your password:")
	pwd = readString()
	return login, pwd
}

// readString read string from command line and remove line-break
func readString() string {
	str := ""
	fmt.Scan(&str)
	return str
}

func (cli *CLI) Start(cancel context.CancelFunc) error {
	if cli.client == nil && !cli.storage.CheckFile() {
		fmt.Println("Can't connect server for the first registration. Sorry, bye!")
		return errors.New("no storage and no server connection")
	}
	if cli.client != nil {
		fmt.Println("-------------- GophKeeper -------------")
	} else {
		fmt.Println("-------- GophKeeper (offline mode) --------")
	}
	for {
		switch cli.state {
		case START:
			q, err := cli.start(cancel)
			if err != nil {
				return err
			}
			if q {
				return nil
			}
			break
		case LOGGED_IN:
			q, err := cli.loggedIn(cancel)
			if err != nil {
				return err
			}
			if q {
				return nil
			}
			break
		}
	}
}

// start unauthorized user interface
func (cli *CLI) start(cancel context.CancelFunc) (quit bool, err error) {
	if cli.client != nil && !cli.storage.CheckFile() {
		fmt.Print("r - registration | ")
	}
	fmt.Println("l - log in | q - exit")
	fmt.Print("-> ")
	command := readString()
	switch command {
	case "r":
		usr, err := cli.registration()
		if err != nil {
			fmt.Println("Registration error, try again later. (" + err.Error() + ") \n")
		} else {
			cli.state = LOGGED_IN
			cli.user = &usr
		}
		break
	case "l":
		usr, err := cli.auth()
		if err != nil {
			fmt.Println("wrong password")
		} else {
			cli.state = LOGGED_IN
			cli.user = &usr
		}
		break
	case "q":
		fmt.Println("Bye!")
		cancel()
		return true, nil
	default:
		fmt.Println("Unknown command " + command)
	}
	return false, nil
}

// addNote get new note's data from user input and add it to storage
func (cli *CLI) addNote() (err error) {
	fmt.Println("Enter note's title")
	title := readString()
	fmt.Println("Enter note's body")
	body := readString()
	err = cli.storage.AddItem(title, body, cli.user.PasswordHash)
	if err != nil {
		return err
	}
	err = cli.storage.Flush(cli.user.PasswordHash)
	return err
}

// loggedIn interface for logged in user
func (cli *CLI) loggedIn(cancel context.CancelFunc) (quit bool, err error) {
	fmt.Println("Welcome, " + cli.user.Login)
	showBar := true
	for {
		syncDate := cli.storage.GetLastSyncDate().Format(time.RFC3339)
		if cli.storage.GetLastSyncDate().IsZero() {
			syncDate = "...never"
		}
		nonSync := cli.storage.GetNonSyncedData()
		if showBar {
			fmt.Printf("You have %d notes. %d notes need sync. Last sync date is %s \n", cli.storage.GetDataLen(), len(nonSync), syncDate)
			if cli.storage.GetDataLen() > 0 {
				fmt.Print("v - view list | v# - view note | ")
			}
			fmt.Print("a - add notes | ")
			if cli.storage.GetDataLen() > 0 {
				fmt.Print("d# - delete note | ")
			}
			if cli.client != nil {
				fmt.Print("s - sync | ")
			}
			fmt.Println("q - exit")
			showBar = false
		}
		fmt.Print("-> ")
		command := readString()
		// delete note
		if matched := deleteCmdRegExp.FindAllStringSubmatch(command, -1); matched != nil {
			index, err := strconv.Atoi(matched[0][1])
			if err != nil {
				fmt.Println("wrong index")
				continue
			}
			err = cli.storage.SetDeleted(index - 1)
			if err != nil {
				fmt.Println("can't delete (" + err.Error() + ")")
				continue
			}
			fmt.Println("note deleted")
			showBar = true
			err = cli.storage.Flush(cli.user.PasswordHash)
			if err != nil {
				fmt.Println("save notes error (" + err.Error() + ")")
			}
			continue
		}
		// view note
		if matched := viewCmdRegExp.FindAllStringSubmatch(command, -1); matched != nil {
			index, _ := strconv.Atoi(matched[0][1])
			n, err := cli.storage.GetByIndex(index - 1)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			encoded, err := crypto.Decrypt(cli.user.PasswordHash, n.Body)
			if err != nil {
				fmt.Println("decrypt error (" + err.Error() + ") \n")
				continue
			}
			fmt.Println(string(encoded))
			continue
		}
		switch command {
		case "a":
			err = cli.addNote()
			if err != nil {
				fmt.Println("add note error (" + err.Error() + ") \n")
			} else {
				fmt.Println("note added!")
				err = cli.storage.Flush(cli.user.PasswordHash)
				showBar = true
				if err != nil {
					fmt.Println("save notes error (" + err.Error() + ")")
				}
			}
		case "v":
			fmt.Println("------------- Notes: -------------")
			loc := time.Local
			for i, note := range cli.storage.GetNotes(true) {
				fmt.Print(strconv.Itoa(i+1) + " | " + note.CreatedAt.In(loc).Format(time.RFC822) + " | " + note.Title)
				if note.Deleted {
					fmt.Print(" (deleted)")
				}
				if note.ID == "" {
					fmt.Print(" (local)")
				}
				fmt.Println("")
			}
			fmt.Println("----------- End notes: -----------")
		case "s":
			newData, newLastSync, err := cli.client.SyncData(nonSync, cli.storage.GetLastSyncDate())
			if err != nil {
				fmt.Println("sync error (" + err.Error() + ")")
				continue
			}
			if cli.storage.GetLastSyncDate().After(newLastSync) {
				fmt.Println("sync error (last sync date is incorrect")
				continue
			}
			err = cli.storage.UpdateData(newData, newLastSync)
			if err != nil {
				fmt.Println("merge data error (" + err.Error() + ")")
				continue
			}
			err = cli.storage.Flush(cli.user.PasswordHash)
			if err != nil {
				fmt.Println("save notes error (" + err.Error() + ")")
			}
			showBar = true
		case "q":
			fmt.Println("Bye!")
			cancel()
			return true, nil
		default:
			fmt.Println("Unknown command " + command)
		}
	}
}

func New(client GRPCClientModel, s StorageModel) *CLI {
	return &CLI{state: START, client: client, storage: s}
}
