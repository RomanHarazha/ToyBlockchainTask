package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
)

func fileCount(path string) (int, error) {
	i := 0
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, err
	}
	for _, file := range files {
		if !file.IsDir() {
			i++
		}
	}
	return i, nil
}

func main() {
	var name, message string
	var enter int
	fmt.Println("Enter:\n1 to create a new blockchain;\n2 to add a block to an existing blockchain;\n3 to validate an existing blockchain;\n0 to exit.")
	for {
		fmt.Print("Enter the command: ")
		fmt.Fscan(os.Stdin, &enter)
		if enter == 0 {
			break
		}
		switch enter {
		case 1:
			fmt.Print("Enter the name of the blockchain you want to create: ")
			fmt.Fscan(os.Stdin, &name)
			err := os.Mkdir(name, 0666)
			if err != nil {
				panic(err)
			}
		case 2:
			fmt.Print("Enter the name of the blockchain you want to add the block to: ")
			fmt.Fscan(os.Stdin, &name)
			path, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			path += "/" + name
			count, err := fileCount(path)
			if err != nil {
				panic(err)
			}
			fmt.Print("Enter the message you want to put in the new block: ")
			fmt.Fscan(os.Stdin, &message)
			file, err := os.Create(path + "/" + fmt.Sprint(count) + ".txt")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			hash := sha256.New()
			hash.Write([]byte(message))
			hashString := hex.EncodeToString(hash.Sum(nil))
			file.WriteString(message + "\r\n" + hashString)
			if count > 0 {
				filePrev, err := os.Open(path + "/" + fmt.Sprint(count-1) + ".txt")
				if err != nil {
					panic(err)
				}
				defer filePrev.Close()
				scanner := bufio.NewScanner(filePrev)
				scanner.Scan()
				scanner.Scan()
				file.WriteString("\r\n" + scanner.Text())
			} else {
				file.WriteString("\r\n" + "0000000000000000000000000000000000000000000000000000000000000000")
			}

		case 3:
			fmt.Print("Enter the name of the blockchain you want to validate: ")
			fmt.Fscan(os.Stdin, &name)
			path, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			path += "/" + name
			count, err := fileCount(path)
			if err != nil {
				panic(err)
			}
			if count > 1 {
				errCatcher := 0
				for i := 0; i < count-1; i++ {
					filePrev, err := os.Open(path + "/" + fmt.Sprint(i) + ".txt")
					if err != nil {
						panic(err)
					}
					defer filePrev.Close()
					ownHash := bufio.NewScanner(filePrev)
					ownHash.Scan()
					ownHash.Scan()
					file, err := os.Open(path + "/" + fmt.Sprint(i+1) + ".txt")
					if err != nil {
						panic(err)
					}
					defer file.Close()
					nextHash := bufio.NewScanner(file)
					nextHash.Scan()
					nextHash.Scan()
					nextHash.Scan()
					if ownHash.Text() != nextHash.Text() {
						errCatcher++
					}
				}
				if errCatcher == 0 {
					fmt.Println("Blockchain is valid.")
				}
			} else if count == 1 {
				fmt.Println("Blockchain consists of a single block!")
			} else {
				fmt.Println("Blockchain is empty!")
			}
		default:
			fmt.Println("There is no such command.")
		}
	}
}
