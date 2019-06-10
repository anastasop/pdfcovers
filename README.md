# pdfcovers

Pdfcovers takes a list of pdf files and generates a pdf containing the first page (cover) of each one.

It can help you cope with archives, usually external disks, which contains pdfs with cryptic names like _TR023213.pdf download.pdf paper.pdf MIT-LCS-1021.pdf_ etc and you no longer remember how you acquired each one. Pdfcovers also adds a text annotation to each page with the full path name of the origin pdf so that you can retrieve it.

## Installation

Pdfcovers is written in Go and tested with go1.12. It can be easily installed with the go tool

```
go get github.com/anastasop/pdfcovers
```

## Usage

```
pdfcovers [-s] [-o file] [-lf <license file>] [-ln <customer name>] [pdf files]
```

The `-o` flag causes it to output the pdf result to file, otherwise to stdout

The `-s` flag causes it to read a lists of filenames from stdin, usually with find like `find . -name *.pdf | pdfcovers -s`. If both arguments files and `-s` are used then the argument files are displayed first in the output. This can be used to add a cover to the collection.

The `-lf` optional flag point to a file containing a license from [unidoc](https://unidoc.io/). See the section on licensing for details.

The `-ln` optional flag is the [unidoc](https://unidoc.io/) customer name for the license. See the section on licensing for details.

## Examples

```
# extract the covers of all pdfs in the current directory
pdfcovers -o index.pdf *.pdf

# extract the covers of all pdfs in the current directory
# and redirect the output to a display program
pdfcovers *.pdf | 9 page

# extract the covers of all pdfs in the current directory tree
find . -name *.pdf | pdfcovers -s -o index.pdf

# extract the covers of all pdfs in the current directory tree
# and add a cover page to the outpus
find . -name *.pdf | pdfcovers -s -o index.pdf cover-page.pdf

# Another way to extract the covers without the -s flag
find . -name *.pdf -exec pdfcovers -o index.pdf '{}' '+'
```

## License

Pdfcovers is licensed under the [AGPL](https://www.gnu.org/licenses/agpl-3.0.en.html) This is a requirement by [unidoc](https://unidoc.io/), the creators of [unipdf](https://github.com/unidoc/unipdf), the library that powers pdfcovers. You can check the above links for details on licensing.

If you don't have a unidoc license you can use pdfcovers but a watermark appears to the bottom of each page. Also it prints `Unlicensed copy of unidoc. To get rid of the watermark - Please get a license on https://unidoc.io`. This is not a major issue as pdfcovers is for personal use and the output is not excepted to be published. If however this is annoying you can [apply to unidoc](https://unidoc.io/pricing/) for a free license and use the `-lf` and `-ln` flag to hide the watermark.

