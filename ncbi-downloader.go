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

func stripTag(input string) (output string) {
	splitString := strings.Split(input, ">")
	//output := ""

	splitString2 := strings.Split(splitString[1], "<")
	output = splitString2[0]
	return
}

func findTag(lines []string, tag string) (output string) {
	for _, line := range lines {
		if strings.Contains(line, tag) {
			output = stripTag(line)
			break
		}
	}
	return
}

func main() {
	fmt.Println("Hello")
	// https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=mopalia+AND+COI&retmax=500
	// https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=nuccore&id=34577062,24475906&rettype=gb&retmode=xml

	id_response, _ := http.Get("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=mopalia+AND+COI&retmax=10")
	htmlData, _ := ioutil.ReadAll(id_response.Body)

	htmlString := string(htmlData)
	splitString := strings.Split(htmlString, "\n")

	for i, line := range splitString {

		if strings.Contains(line, "<Id>") {
			fmt.Println(i, " => ", stripTag(line))
		}
	}

	xml, err := ioutil.ReadFile("sequence.gbx.xml")
	if err != nil {
		fmt.Print(err)
	}
	xmlString := string(xml)
	xmlLines := strings.Split(xmlString, "\n")

	locus_id := findTag(xmlLines, "GBSeq_locus")
	fmt.Println(locus_id)

	//fmt.Println(xmlString)
	//fmt.Println(reflect.TypeOf(xmlString))

	//f, _ := os.Create("sequence.gb")
	//check(err)
	//defer f.Close()
	//test := "abcd"
	//f.WriteString(htmlString)
	//f.Sync()
}
