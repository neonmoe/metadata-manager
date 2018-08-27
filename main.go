package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

func main() {
	var fs = &fasthttp.FS{
		Root:       "./public",
		IndexNames: []string{"index.html"},
	}
	var fsHandler = fs.NewRequestHandler()

	// Load the index file for cutting and slicing
	var indexHtmlRaw, err = ioutil.ReadFile("./public/index.html")
	if err != nil {
		panic(err)
	}
	var indexHtml = string(indexHtmlRaw)

	// Save the pre- and post-custom-field parts of the index file
	var indexHtmlParts = [2]string{}
	copy(indexHtmlParts[:], strings.Split(indexHtml, "<!-- CUSTOM FIELDS -->")[0:2])
	// Extract the template from the index file
	var templateField = strings.Split(strings.Split(indexHtml, "<!--FIELD TEMPLATE")[1],
		"FIELD TEMPLATE-->")[0]
	// Pre-render with the default record
	var defaultRecord = [][]string{
		[]string{"Title", "Creator", "Identifier", "Dates", "Subject", "Funders",
			"Rights", "Language", "Location", "Methodology"},
		[]string{"", "", "", "", "", "", "", "", "", "", ""},
	}
	var defaultIndexHtml bytes.Buffer
	writeIndexHtml(&defaultIndexHtml, indexHtmlParts, templateField, defaultRecord)

	// Finally, create the server
	var server = &Server{
		fsHandler:        fsHandler,
		indexHtmlParts:   indexHtmlParts,
		defaultIndexHtml: defaultIndexHtml.Bytes(),
		templateField:    templateField,
	}

	var addr = "0.0.0.0:8080"
	fmt.Printf("Serving at: %s", addr)
	fasthttp.ListenAndServe(addr, server.HandleRequest)
}

type Server struct {
	fsHandler        fasthttp.RequestHandler
	indexHtmlParts   [2]string
	defaultIndexHtml []byte
	templateField    string
}

func (s *Server) HandleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("text/html")
		var returnDefaultPage = true

		// Handle imported files (return a page with the imported data set)
		var importFile, fileErr = ctx.FormFile("import-file")
		if fileErr == nil {
			var file, err = importFile.Open()
			if err == nil {
				var formRecord = createRecordFromCsv(file)
				if len(formRecord) > 0 {
					writeIndexHtml(ctx, s.indexHtmlParts, s.templateField, formRecord)
					returnDefaultPage = false
				}
			}
		}

		// Handle new fields being added (return a page with the new field)
		var queryString = ctx.PostArgs().QueryString()
		if returnDefaultPage && len(queryString) > 0 {
			var formRecord, newField = createRecordFromQueryString(string(ctx.PostArgs().QueryString()))
			formRecord[0] = append(formRecord[0], newField[0])
			formRecord[1] = append(formRecord[1], newField[1])
			writeIndexHtml(ctx, s.indexHtmlParts, s.templateField, formRecord)
			returnDefaultPage = false
		}

		if returnDefaultPage {
			ctx.Write(s.defaultIndexHtml)
		}

	case "/export":
		var formRecord, _ = createRecordFromQueryString(string(ctx.PostArgs().QueryString()))
		var id = getValueFromRecord(formRecord, "Identifier")
		var name string
		if len(id) > 0 {
			name = fmt.Sprintf("attachment; filename=\"%s-metadata.csv\"", id)
		} else {
			name = "attachment; filename=\"metadata.csv\""
		}
		ctx.Response.Header.Set("Content-Disposition", name)
		var csvWriter = csv.NewWriter(ctx)
		csvWriter.WriteAll(formRecord)
	default:
		s.fsHandler(ctx)
	}
}

func writeIndexHtml(writer io.Writer, bodyParts [2]string, templateField string, fieldRecord [][]string) {
	writer.Write([]byte(bodyParts[0]))
	for i := range fieldRecord[0] {
		if len(fieldRecord[0][i]) > 0 {
			writeNewFieldHtml(writer, templateField, fieldRecord[0][i], fieldRecord[1][i])
		}
	}
	writer.Write([]byte(bodyParts[1]))
}

func writeNewFieldHtml(writer io.Writer, templateField string, title string, value string) {
	var replacer = strings.NewReplacer("TITLE", title, "VALUE", value)
	writer.Write([]byte(replacer.Replace(templateField)))
}

func createRecordFromQueryString(queryString string) ([][]string, [2]string) {
	var formParams = strings.Split(queryString, "&")
	var newField = [2]string{}
	var records = make([][]string, 2)
	for i := range records {
		records[i] = make([]string, len(formParams))
	}
	for i, param := range formParams {
		var paramParts = strings.Split(param, "=")
		var hasTwoParts = len(paramParts) == 2
		if paramParts[0] == "AddedFieldName" {
			if hasTwoParts {
				newField[0] = tryToQueryUnescape(paramParts[1])
			}
		} else if paramParts[0] == "AddedFieldValue" {
			if hasTwoParts {
				newField[1] = tryToQueryUnescape(paramParts[1])
			}
		} else {
			records[0][i] = tryToQueryUnescape(paramParts[0])
			if hasTwoParts {
				records[1][i] = tryToQueryUnescape(paramParts[1])
			}
		}
	}
	return records, newField
}

func createRecordFromCsv(reader io.Reader) [][]string {
	var csvReader = csv.NewReader(reader)
	var records, err = csvReader.ReadAll()
	if err != nil {
		return make([][]string, 0)
	} else {
		return records
	}
}

func getValueFromRecord(formRecord [][]string, key string) string {
	for i, k := range formRecord[0] {
		if k == key {
			return formRecord[1][i]
		}
	}
	return ""
}

func tryToQueryUnescape(s string) string {
	var unescaped, err = url.QueryUnescape(s)
	if err == nil {
		return unescaped
	} else {
		return s
	}
}
