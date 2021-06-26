package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s []int
}

func NewStack() *stack {
	return &stack {sync.Mutex{}, make([]int,0), }
}

func (s *stack) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()


	l := len(s.s)
	if l == 0 {
		return 0, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func getOpcode(instruction string) int {
	opcodes := map[string] int {
		"wakeup": 0,
		"speak": 1,
		"sleep": 2,
		"push": 3,
		"pop": 4,
		"destruct": 5,
	}

	opcode, ok := opcodes[instruction]

	if ok == false {
		return -1
	} else {
		return opcode
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {
	//making this language without knowing 1 bit of go

	if os.Args[1] == "" {
		fmt.Println("provide the file name with your code")
		os.Exit(404)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	s := string(data)

	//encryptFile("test.txt", data, string(rand2.Int()))

	i := strings.Fields(s)

	strings := []string{}

	stack1 := NewStack()

	printCount := 0
	pushCount := 0

	inString := false

	currentString := ""

	if string(s[len(s)-len("sleep"):]) != "sleep" {
		log.Fatal(errors.New("the turtle fucking died from being awake for too long"))
	}

	if string(s[0:len("wakeup")]) != "wakeup" {
		log.Fatal(errors.New("zzzzzzzzzzzzzzzzzzzzzzzzzz"))
	}

	for _, element := range s {
		//fmt.Println(string(element))

		if string(element) == "'" && inString == true {
			inString = false
			strings = append(strings, currentString)
			//fmt.Println(currentString)
			currentString = ""

			continue
		}

		if string(element) == "'" && inString == false {
			inString = true
			continue
		}

		if inString == true {
			currentString += string(element)

			continue
		}
	}

	for _, element := range i {
		//fmt.Println(getOpcode(element))
		//fmt.Println(element)
		switch getOpcode(element) {
		case 0:
			fmt.Println("the turtle is waking up...")
			time.Sleep(2 * time.Second)
		case 1:
			fmt.Println(strings[printCount])

			printCount += 1
			pushCount += 1
		case 2:
			fmt.Println("going to sleep...")
			time.Sleep(2 * time.Second)
			fmt.Println("the turtle is asleep")

			os.Exit(30000)
		case 3:
			for _, char := range strings[pushCount] {
				stack1.Push(int(char))
			}

			pushCount += 1
			printCount += 1
		case 4:
			poppedValue, _ := stack1.Pop()

			fmt.Println(string(poppedValue))
		case 5:
			h := sha256.New()
			h.Write([]byte(os.Args[1]))

			if err := os.WriteFile(os.Args[1], []byte(h.Sum(nil)), 0666); err != nil {
				log.Fatal(err)
			}

			fmt.Println("💥 destroyed")
		default:

		}
	}
}