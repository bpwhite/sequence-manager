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

	// Returns value from XML line.

	splitString := strings.Split(input, ">")
	//output := ""

	splitString2 := strings.Split(splitString[1], "<")
	output = splitString2[0]
	return
}


func findTag(lines []string, tag string, offset int) (output string) {

	// Finds first occurrence of <tag> in lines (XML) and returns value
	// from line that is offset number of lines beyond <tag>.

	output = "NA"

	for i, line := range lines {
		if strings.Contains(line, tag) {
			output = strings.Replace(stripTag(lines[i+offset]), ",", "_", -1)
			break
		}
	}
	return
}


func findTags(lines []string, tag string) (output string) {

	// Finds first occurrence of <tag> in lines (XML format) and returns
	// semicolon-concatenated values for all sub-tags. Written with 
	// <GBReference_authors> in mind for which there are multiple <GBAuthor> subtags

	end_tag := strings.Replace(tag, "<", "</", 1)
	collect := false

	for _, line := range lines {

		if collect == true {

			val := stripTag(line)

			if output != "" {

				output = fmt.Sprintf("%s;%s", output, val)

			} else {

				output = val

			}

		}

		if strings.Contains(line, end_tag) {

			break

		} else if strings.Contains(line, tag) {

			collect = true

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
	t := time.Now()
	//fmt.Println(t.Format("20060102150405"))
	time_stamp := t.Format("20060102150405")
	status_string := fmt.Sprint("output/", time_stamp, "_", run_tag, "_", *taxonPtr, ".html")
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


	// Count how many records were found
	seq_counter1 := 0 
	for _, line := range splitString {

		if strings.Contains(line, "<Id>") {
			outht.WriteString(fmt.Sprint(line, ","))
			seq_counter1++
		}
	}
	seq_counter2 := 0


	// write headers for fields
	fmt.Println(	`locus_id,seq_length,strandedness,moltype,toplogy,division,update_date,create_date,definition,primary_accession,accession_version,source,organism,taxonomy,nuc_sequence,prot_sequence,taxon_id,gene,product,codon_start,organelle,pub_title,pub_authors,pub_jrn,voucher,country,lat_long,note`)

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

			// Get pertinent information from XML
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
			GB_taxon_id := strings.SplitAfter(findTag(xmlLines, "<GBQualifier_name>db_xref", 1), "taxon:")[1]
			GB_gene := findTag(xmlLines, "<GBQualifier_name>gene", 1)
			GB_product := findTag(xmlLines, "<GBQualifier_name>product", 1)
			GB_codon_start := findTag(xmlLines, "<GBQualifier_name>codon_start", 1)
			GB_organelle := findTag(xmlLines, "<GBQualifier_name>organelle", 1)
			GB_pub_title := findTag(xmlLines, "<GBReference_title>", 0)
			GB_pub_jrn := findTag(xmlLines, "<GBReference_journal>", 0)
			GB_pub_authors := strings.TrimRight(strings.Replace(findTags(xmlLines, "<GBReference_authors>"), ",", " ", -1), ";") // multiple tags

			GB_voucher := findTag(xmlLines, "<GBQualifier_name>specimen_voucher", 1)
			GB_country := findTag(xmlLines, "<GBQualifier_name>country", 1)
			GB_lat_long := findTag(xmlLines, "<GBQualifier_name>lat_long", 1)
			GB_note := findTag(xmlLines, "<GBQualifier_name>note", 1)


			// Print CSV lines to stdout
			outputFields := []string {GB_locus_id,
						GB_seq_length,
						GB_strandedness,
						GB_moltype,
						GB_toplogy,
						GB_division,
						GB_update_date,
						GB_create_date,
						GB_definition,
						GB_primary_accession,
						GB_accession_version,
						GB_source,
						GB_organism,
						GB_taxonomy,
						GB_nuc_sequence,
						GB_prot_sequence,
						GB_taxon_id,
						GB_gene,
						GB_product,
						GB_codon_start,
						GB_organelle,
						GB_pub_title,
						GB_pub_authors,
						GB_pub_jrn,
						GB_voucher,
						GB_country,
						GB_lat_long,
						GB_note}

			fmt.Println(strings.Join(outputFields,","))

			// Make sure not to exceed reported ID's
			if seq_counter2 == seq_counter1 {
				break
			}

		}
	}

}
