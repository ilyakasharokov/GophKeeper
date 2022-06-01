// cli - command line interface of the app
package cli

import (
	"context"
	"errors"
	"fmt"
	"gophkeeper/pkg/models"

	"gophkeeper/pkg/crypto"
	"regexp"
	"strconv"
	"time"
)

const (
	startState = iota
	loggedInState
)

type CLI struct {
	state   int
	auther  Authenticator
	storage NotesKeeper
	syncer  Syncer
	user    *models.User
}

type NotesKeeper interface {
	AddNote(title string, body string, key []byte) error
	GetNotes(all bool) []models.Note
	GetNote(index int) (models.Note, error)
	DeleteNote(index int) error
	GetNotesCount() int
	Flush(hash []byte) error
	CheckFile() bool
	Load(hash []byte) error
}

type Syncer interface {
	Sync() error
	GetNonSyncNotes() []models.Note
	GetLastSyncDate() time.Time
	UpdateLastSyncDate()
}

type Authenticator interface {
	Login(login string, pwd string) (string, error)
	Registration(login string, pwd string) (string, string, error)
}

var deleteCmdRegExp = regexp.MustCompile(`^d(\d+)$`)
var viewCmdRegExp = regexp.MustCompile(`^v(\d+)$`)

// registration register and return user if success
func (cli *CLI) registration() (models.User, error) {
	login, pwd := readAuth()
	status, token, err := cli.auther.Registration(login, pwd)
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
	if cli.auther != nil { // login remote
		token, err = cli.auther.Login(login, pwd)
		if err != nil {
			return models.User{}, errors.New("can't auth remotely")
		}
	}
	if cli.storage.CheckFile() { // local storage file found
		err = cli.storage.Load(hash)
		if err != nil {
			return models.User{}, err
		}
		cli.syncer.UpdateLastSyncDate()
	} else { // new client for existing user
		err := cli.syncer.Sync()
		if err != nil {
			return models.User{}, err
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
	_, err := fmt.Scan(&str)
	if err != nil {
		fmt.Println(err.Error())
	}
	return str
}

// Start CLI interface
func (cli *CLI) Start(cancel context.CancelFunc) error {
	if cli.auther == nil && !cli.storage.CheckFile() {
		fmt.Println("Can't connect server for the first registration. Sorry, bye!")
		return errors.New("no storage and no server connection")
	}
	if cli.auther != nil {
		fmt.Println("-------------- GophKeeper -------------")
	} else {
		fmt.Println("-------- GophKeeper (offline mode) --------")
	}
	for {
		switch cli.state {
		case startState:
			q, err := cli.start(cancel)
			if err != nil {
				return err
			}
			if q {
				return nil
			}
			break
		case loggedInState:
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
	if cli.auther != nil && !cli.storage.CheckFile() {
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
			cli.state = loggedInState
			cli.user = &usr
		}
		break
	case "l":
		usr, err := cli.auth()
		if err != nil {
			fmt.Println("wrong password")
		} else {
			cli.state = loggedInState
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
func (cli *CLI) addNote() error {
	fmt.Println("Enter note's title")
	title := readString()
	fmt.Println("Enter note's body")
	body := readString()
	err := cli.storage.AddNote(title, body, cli.user.PasswordHash)
	if err != nil {
		fmt.Println("add note error (" + err.Error() + ") \n")
	} else {
		fmt.Println("note added!")
		err = cli.storage.Flush(cli.user.PasswordHash)
		if err != nil {
			fmt.Println("save storage file error (" + err.Error() + ")")
		}
	}
	return nil
}

// showMenuBar output menu bar
func (cli *CLI) showMenuBar(notesCount int, notesToSyncCount int, lastSyncDate string) {
	fmt.Printf("You have %d notes. %d notes need sync. Last sync date is %s \n", notesCount, notesToSyncCount, lastSyncDate)
	fmt.Print("v - view list | v# - view note | a - add notes | d# - delete note |")
	if cli.auther != nil {
		fmt.Print("s - sync | ")
	}
	fmt.Println("q - exit")
}

// deleteNote delete note from storage
func (cli *CLI) deleteNote(inputIndex string) error {
	index, err := strconv.Atoi(inputIndex)
	if err != nil {
		return errors.New("wrong index")
	}
	err = cli.storage.DeleteNote(index - 1)
	if err != nil {
		return errors.New("can't delete (" + err.Error() + ")")
	}
	err = cli.storage.Flush(cli.user.PasswordHash)
	if err != nil {
		fmt.Println("save storage file error (" + err.Error() + ")")
	}
	return nil
}

// viewNote output note's content
func (cli *CLI) viewNote(inputIndex string) error {
	viewIndex, _ := strconv.Atoi(inputIndex)
	n, err := cli.storage.GetNote(viewIndex - 1)
	if err != nil {
		return err
	}
	body := n.Body
	encoded, err := crypto.Decrypt(cli.user.PasswordHash, body)
	if err != nil {
		return errors.New("decrypt error (" + err.Error() + ") \n")
	}
	fmt.Println(string(encoded))
	return nil
}

// loggedIn interface for logged in user
func (cli *CLI) loggedIn(cancel context.CancelFunc) (quit bool, err error) {
	fmt.Println("Welcome, " + cli.user.Login)
	showBar := true
	for {
		syncDate := cli.syncer.GetLastSyncDate().Format(time.RFC3339)
		if cli.syncer.GetLastSyncDate().IsZero() {
			syncDate = "...never"
		}
		nonSync := cli.syncer.GetNonSyncNotes()
		if showBar {
			cli.showMenuBar(cli.storage.GetNotesCount(), len(nonSync), syncDate)
			showBar = false
		}
		fmt.Print("-> ")
		command := readString()
		// delete note
		if matched := deleteCmdRegExp.FindAllStringSubmatch(command, -1); matched != nil {
			err := cli.deleteNote(matched[0][1])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				showBar = true
				fmt.Println("note deleted")
			}
			continue
		}
		// view note
		if matched := viewCmdRegExp.FindAllStringSubmatch(command, -1); matched != nil {
			err := cli.viewNote(matched[0][1])
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		switch command {
		case "a":
			cli.addNote()
		case "v":
			cli.printNotesTable()
		case "s":
			cli.sync()
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

// sync data
func (cli *CLI) sync() {
	err := cli.syncer.Sync()
	if err != nil {
		fmt.Println(err)
	} else {
		err := cli.storage.Flush(cli.user.PasswordHash)
		if err != nil {
			fmt.Println("save storage file error (" + err.Error() + ")")
		} else {
			fmt.Println("synchronization was successfull!")
		}

	}
}

// printNotesTable output notes info table
func (cli *CLI) printNotesTable() {
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
	fmt.Println("----------------------")
}

func New(auther Authenticator, s NotesKeeper, syncer Syncer) *CLI {
	return &CLI{state: startState, auther: auther, storage: s, syncer: syncer}
}
