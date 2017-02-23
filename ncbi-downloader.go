package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
	"net/http"
	//"os"
	//"reflect"
	"flag"
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

func findTag(lines []string, tag string, offset int) (output string) {
	for i, line := range lines {
		if strings.Contains(line, tag) {
			output = stripTag(lines[i+offset])
			break
		}
	}
	return
}

func main() {
	fmt.Println("Hello")
	// https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=mopalia+AND+COI&retmax=500
	// https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=nuccore&id=34577062,24475906&rettype=gb&retmode=xml

	// Parse command line arguments
	stermPtr := flag.String("sterm", "COI", "a string")
	taxonPtr := flag.String("taxon", "mopalia", "a string")
	retmaxPtr := flag.Int("retmax", 10, "an int")
	//boolPtr := flag.Bool("test", false, "a bool")

	flag.Parse()

	fmt.Println("search term:", *stermPtr)
	fmt.Println("taxon:", *taxonPtr)
	fmt.Println("retmax:", *retmaxPtr)
	// End command line arguments

	// Concatenate esearch string
	concat_string := fmt.Sprint("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&term=", *taxonPtr, "+AND+", *stermPtr, "&retmax=", *retmaxPtr, "")
	fmt.Println(concat_string, "\n")
	id_response, _ := http.Get(concat_string)
	htmlData, _ := ioutil.ReadAll(id_response.Body)

	htmlString := string(htmlData)
	splitString := strings.Split(htmlString, "\n")

	for i, line := range splitString {

		if strings.Contains(line, "<Id>") {
			fmt.Println(i, " => ", stripTag(line))
			gb_id := stripTag(line)
			concat_request := fmt.Sprint("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=nuccore&id=", gb_id, "&rettype=gb&retmode=xml")
			fmt.Println("Requesting...", concat_request)
			gb_response, _ := http.Get(concat_request)
			gb_data, _ := ioutil.ReadAll(gb_response.Body)

			xmlString := string(gb_data)
			xmlLines := strings.Split(xmlString, "\n")

			GB_locus_id := findTag(xmlLines, "GBSeq_locus", 0)
			GB_seq_length := findTag(xmlLines, "GBSeq_length", 0)
			GB_strandedness := findTag(xmlLines, "GBSeq_strandedness", 0)
			GB_moltype := findTag(xmlLines, "GBSeq_moltype", 0)
			GB_toplogy := findTag(xmlLines, "GBSeq_topology", 0)
			GB_division := findTag(xmlLines, "GBSeq_division", 0)
			GB_update_date := findTag(xmlLines, "GBSeq_update-date", 0)
			GB_create_date := findTag(xmlLines, "GBSeq_create-date", 0)
			GB_definition := findTag(xmlLines, "GBSeq_definition", 0)
			GB_primary_accession := findTag(xmlLines, "GBSeq_primary-accession", 0)
			GB_accession_version := findTag(xmlLines, "GBSeq_accession-version", 0)
			GB_nuc_sequence := findTag(xmlLines, "GBSeq_sequence", 0)
			GB_source := findTag(xmlLines, "GBSeq_source", 0)
			GB_organism := findTag(xmlLines, "GBSeq_organism", 0)
			GB_taxonomy := findTag(xmlLines, "GBSeq_taxonomy", 0)
			GB_prot_sequence := findTag(xmlLines, "<GBQualifier_name>translation", 1)

			fmt.Println(GB_locus_id)
			fmt.Println(GB_seq_length)
			fmt.Println(GB_strandedness)
			fmt.Println(GB_moltype)
			fmt.Println(GB_toplogy)
			fmt.Println(GB_division)

			fmt.Println(GB_update_date)
			fmt.Println(GB_create_date)
			fmt.Println(GB_strandedness)
			fmt.Println(GB_definition)
			fmt.Println(GB_primary_accession)
			fmt.Println(GB_accession_version)
			fmt.Println(GB_source)
			fmt.Println(GB_organism)
			fmt.Println(GB_taxonomy)

			fmt.Println(GB_nuc_sequence)
			fmt.Println(GB_prot_sequence)
		}
	}
	//fmt.Println(xmlString)
	//fmt.Println(reflect.TypeOf(xmlString))

	//f, _ := os.Create("sequence.gb")
	//check(err)
	//defer f.Close()
	//test := "abcd"
	//f.WriteString(htmlString)
	//f.Sync()
}
