package gokeeper

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"terminal"
)

var (
	STORAGE_PATH string = ".gokeeper.db"
	KEY          []byte
)

func GetPass(msg string) string {
	fmt.Printf(msg)
	key, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Printf("\n")
	if err != nil {
		return ""
	}
	return string(key)
}

func GetInput(msg string) string {
	fmt.Printf(msg)
	res, _, err := bufio.NewReader(os.Stdin).ReadLine()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(res)
}

func add(s *Storage) error {
	key := GetInput("Key : ")
	data := []byte(GetInput("Data : "))
	err := s.Put(key, data, KEY)
	//s.Save()
	return err
}

func list(s *Storage) {
	var keys []string
	for key, _ := range s.Data() {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Println(key)
	}
}

func del(s *Storage) {
	key := GetInput("Key : ")
	delete(s.Data(), key)
	//s.Save()
}

func show(s *Storage, key string) error {
	data, err := s.Get(key, KEY)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func update_key(s *Storage) error {
	newmasterkey := Skein256([]byte(GetInput("New Master Key : ")))
	err := s.UpdateKey(KEY, newmasterkey)
	if err != nil {
		return err
	}
	KEY = newmasterkey
	return nil
}

func Main() {
	if len(os.Args) > 1 {
		STORAGE_PATH = os.Args[1]
	}
	KEY = Skein256([]byte(GetPass("Master key : ")))
	storage, err := NewStorage(STORAGE_PATH, KEY)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println(err)
		return
	}
	var command string = ""
	for {
		command = GetInput("> ")
		switch command {
		case "add", "a":
			err := add(storage)
			if err != nil {
				fmt.Println(err)
			}
		case "list", "l":
			list(storage)
		case "del", "d":
			del(storage)
		case "save", "s":
			err := storage.Save(KEY)
			if err != nil {
				fmt.Println(err)
			}
		case "update", "u":
			err := update_key(storage)
			if err != nil {
				fmt.Println(err)
			}
		case "quit", "q":
			return
		default:
			err := show(storage, command)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
