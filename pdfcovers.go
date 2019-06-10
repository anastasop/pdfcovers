package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	pdflicense "github.com/unidoc/unipdf/v3/common/license"
	pdfcore "github.com/unidoc/unipdf/v3/core"
	pdf "github.com/unidoc/unipdf/v3/model"
)

var usageMessage = `usage: pdfcovers [-s] [-o file] [pdf files]

Reads pdfs and creates a new pdf containing the first page of each.
If a pdf is corrupted or encrypted and cannot extract the cover,
then it logs it and proceeds with the next one.

In every page it adds a text annotation with the filename of the origin pdf.

Examples:
  pdfcovers -o index.pdf *.pdf
  pdfcovers *.pdf | 9 page
  find . -name *.pdf | pdfcovers -s -o index.pdf
  find . -name *.pdf | pdfcovers -s -o index.pdf cover-page.pdf
  find . -name *.pdf -exec pdfcovers -s -o index.pdf '{}' '+'

Flags:
`

func usage() {
	fmt.Fprintf(os.Stderr, usageMessage)
	flag.PrintDefaults()
	os.Exit(2)
}

var outputFile = flag.String("o", "", "The output filename. The default is stdout")
var sourceStdin = flag.Bool("s", false, "Read file names from stdin. Useful with 'find | pdfcovers -s'")
var licenseFile = flag.String("lf", "", "A file with a unidoc(http://unidoc.io/pricing) license")
var customerName = flag.String("ln", "", "The name of the customer with a unidoc(http://unidoc.io/pricing) license")

func setLicense() {
	if *licenseFile == "" && *customerName == "" {
		return
	}

	licenseKey, err := ioutil.ReadFile(*licenseFile)
	if err != nil {
		log.Fatal("Cannot read the unidoc license file: ", err)
	}

	if err := pdflicense.SetLicenseKey(string(licenseKey), *customerName); err != nil {
		log.Fatal("Cannot set the unidoc license: ", err)
	}
}

func appendPage(w *pdf.PdfWriter, fname string, pageNum int) error {
	fin, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fin.Close()

	r, err := pdf.NewPdfReader(fin)
	if err != nil {
		return err
	}

	page, err := r.GetPage(pageNum)
	if err != nil {
		return err
	}

	ann := pdf.NewPdfAnnotationText()
	ann.Contents = pdfcore.MakeString(fname)
	ann.Rect = pdfcore.MakeArrayFromIntegers([]int{20, 100, 60, 150})
	page.AddAnnotation(ann.PdfAnnotation)

	err = w.AddPage(page)
	if err != nil {
		return err
	}

	return nil
}

func writeOutput(w *pdf.PdfWriter) error {
	var fout io.WriteSeeker

	if *outputFile == "" {
		fout = os.Stdout
	} else {
		f, err := os.Create(*outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		fout = f
	}

	return w.Write(fout)
}

func main() {
	log.SetPrefix("")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() == 0 && !*sourceStdin {
		usage()
	}

	setLicense()

	w := pdf.NewPdfWriter()

	if flag.NArg() > 0 {
		for _, fname := range flag.Args() {
			if err := appendPage(&w, fname, 1); err != nil {
				log.Println("Failed to add pages from "+fname+":", err)
			}
		}
	}

	if *sourceStdin {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fname := scanner.Text()
			if err := appendPage(&w, fname, 1); err != nil {
				log.Println("Failed to add pages from "+fname+":", err)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println("Failed to read filenames from stdin:", err)
		}
	}

	if err := writeOutput(&w); err != nil {
		log.Println("Failed to write output pdf:", err)
	}
}
