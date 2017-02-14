package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
	"net/http"
	//"os"
	//"reflect"
	"strings"
	//"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("Hello")
	// https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=mopalia+AND+COI&retmax=500

	response, _ := http.Get("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=mopalia+AND+COI&retmax=10")
	/*
		if err != nil {
			log.Fatal(err)
		} else {
			defer response.Body.Close()
			_, err := io.Copy(os.Stdout, response.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
	*/

	//words := strings.Fields(response)
	//fmt.Println(response)

	htmlData, _ := ioutil.ReadAll(response.Body)
	/*
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	*/
	// print out
	htmlString := string(htmlData)
	splitString := strings.Split(htmlString, "\n")
	fmt.Println(len(splitString))

	for i, line := range splitString {
		fmt.Println(i, " => ", line)
	}
	//fmt.Println(reflect.TypeOf(response))

	//f, _ := os.Create("temp.txt")
	//check(err)
	//defer f.Close()
	//test := "abcd"
	//f.WriteString(htmlString)
	//f.Sync()
}
