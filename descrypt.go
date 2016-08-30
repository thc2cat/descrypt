package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/nyarla/go-crypt"
)

var wg sync.WaitGroup
var hashmap = make(map[string]bool)
var dict = make(map[string]bool)

// load hasfile,dictfile in maps,
// launch concurent DES crypts segmented around hash salts
// Output : Found xxdictwords in yyyhashword
func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage : ", os.Args[0], " hashfile dictfile")
	}
	hashmap = readFileIntoMapofString(os.Args[1])
	dict = readFileIntoMapofString(os.Args[2])
	concurentcrypt()
}

// cryptsalt is a crypt.Crypt  loop using keydict map  for a given salt
// synchro with global wg syncGroup for concurency
func cryptsalt(salt string) {

	defer wg.Done()
	for keydict := range dict {
		str := crypt.Crypt(keydict, salt)
		// fmt.Println("key : ", key, keyhash[0:2], str)
		if hashmap[str] {
			fmt.Println("Found :", keydict, "\tin\t", str)
		}
	}
}

// concurentcrypt loop aroung hasmap salts and cryptsalt func
// for concurent crypt using wg syncGroup
func concurentcrypt() {
	salts := make(map[string]bool)
	// Crypt("testtest", "es") =>  `esDRYJnY4VaGM`
	// fmt.Println( "%T(%v)\n" , hashmap )
	// salts <= hashmap
	for keyhash := range hashmap {
		salts[keyhash[0:2]] = true
	}
	fmt.Println(len(salts), " salts, ", len(dict), "words")
	for salt := range salts {
		//go fmt.Println(" cryptsalt(dict,hashmap,", salt, ") ")
		wg.Add(1)
		go cryptsalt(salt)
	}
	wg.Wait()
}

// Small check of err eroor , used to stop quickly
func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

// read filename strings content into a map[string]bool
func readFileIntoMapofString(filename string) map[string]bool {
	myhashmap := make(map[string]bool)
	// hash type is :  []byte
	hash, err := ioutil.ReadFile(filename)
	checkErr(err)
	for _, s := range bytes.Split(hash, []byte("\n")) {
		str := string(s)
		if len(str) > 0 && str != "" {
			myhashmap[str] = true
		}
	}
	return myhashmap
}
