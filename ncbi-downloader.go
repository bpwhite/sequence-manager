package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"time"
	//"reflect"
	"flag"
	"strings"
	//"github.com/davecgh/go-spew/spew"
	"math/rand"
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
			output = strings.Replace(stripTag(lines[i+offset]), ",", "_", -1)
			break
		}
	}
	return
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {

	// Parse command line arguments
	stermPtr := flag.String("sterm", "COI", "a string")
	taxonPtr := flag.String("taxon", "mopalia", "a string")
	retmaxPtr := flag.Int("retmax", 1, "an int")
	//boolPtr := flag.Bool("test", false, "a bool")

	flag.Parse()

	// output
	run_tag := RandStringRunes(6)
	output_string := fmt.Sprint("output/", run_tag, "_", *taxonPtr, ".csv")
	outp, _ := os.Create(output_string)
	status_string := fmt.Sprint("output/", run_tag, "_", *taxonPtr, ".html")
	outht, _ := os.Create(status_string)

	outht.WriteString(fmt.Sprint("search term:", *stermPtr, "<br>"))
	outht.WriteString(fmt.Sprint("taxon:", *taxonPtr, "<br>"))
	outht.WriteString(fmt.Sprint("retmax:", *retmaxPtr, "<br>"))
	// End command line arguments

	// Concatenate esearch string
	concat_string := fmt.Sprint("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=nucleotide&retmax=", *retmaxPtr, "&term=(", *taxonPtr, "+AND+(", *stermPtr, "))")
	outht.WriteString(fmt.Sprint(concat_string, "<br>"))
	//os.Exit(1)
	id_response, _ := http.Get(concat_string)
	htmlData, _ := ioutil.ReadAll(id_response.Body)

	htmlString := string(htmlData)
	splitString := strings.Split(htmlString, "\n")

	seq_counter1 := 0
	for _, line := range splitString {

		if strings.Contains(line, "<Id>") {
			outht.WriteString(fmt.Sprint(line, ","))
			seq_counter1++
		}
	}
	seq_counter2 := 0

	//check(err)
	//defer outp.Close()
	for _, line := range splitString {

		if strings.Contains(line, "<Id>") {
			seq_counter2++
			outht.WriteString(fmt.Sprint(line, "<br>"))
			outht.WriteString(fmt.Sprint(seq_counter2, "/", seq_counter1, " => ", stripTag(line), "<br>"))
			gb_id := stripTag(line)
			concat_request := fmt.Sprint("https://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=nuccore&id=", gb_id, "&rettype=gb&retmode=xml&retmax=1")
			outht.WriteString(fmt.Sprint("Requesting...", concat_request, "<br>"))
			// Sleep between html requests
			time.Sleep(time.Millisecond * 500)
			// Request html page
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
			GB_organism := strings.Replace(findTag(xmlLines, "GBSeq_organism", 0), " ", "_", -1)
			GB_taxonomy := findTag(xmlLines, "GBSeq_taxonomy", 0)
			GB_prot_sequence := findTag(xmlLines, "<GBQualifier_name>translation", 1)

			outp.WriteString(fmt.Sprint(GB_locus_id, ","))
			outp.WriteString(fmt.Sprint(GB_seq_length, ","))
			outp.WriteString(fmt.Sprint(GB_strandedness, ","))
			outp.WriteString(fmt.Sprint(GB_moltype, ","))
			outp.WriteString(fmt.Sprint(GB_toplogy, ","))
			outp.WriteString(fmt.Sprint(GB_division, ","))

			outp.WriteString(fmt.Sprint(GB_update_date, ","))
			outp.WriteString(fmt.Sprint(GB_create_date, ","))
			outp.WriteString(fmt.Sprint(GB_definition, ","))
			outp.WriteString(fmt.Sprint(GB_primary_accession, ","))
			outp.WriteString(fmt.Sprint(GB_accession_version, ","))
			outp.WriteString(fmt.Sprint(GB_source, ","))
			outp.WriteString(fmt.Sprint(GB_organism, ","))
			outp.WriteString(fmt.Sprint(GB_taxonomy, ","))

			outp.WriteString(fmt.Sprint(GB_nuc_sequence, ","))
			outp.WriteString(fmt.Sprint(GB_organism, "$", GB_nuc_sequence, ","))
			outp.WriteString(fmt.Sprint(GB_prot_sequence, ","))
			outp.WriteString("\n")
			// Make sure not to exceed reported ID's
			if seq_counter2 == seq_counter1 {
				break
			}

		}
	}

	outp.Sync()
	//fmt.Println(xmlString)
	//fmt.Println(reflect.TypeOf(xmlString))

	//f, _ := os.Create("sequence.gb")
	//check(err)
	//defer f.Close()
	//test := "abcd"
	//f.WriteString(htmlString)
	//f.Sync()
}
