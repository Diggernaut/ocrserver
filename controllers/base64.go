package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract"
	"github.com/otiai10/marmoset"
)

// Base64 ...
func Base64(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	render := marmoset.Render(w, true)

	var body = new(struct {
		Base64    string                `json:"base64"`
		Trim      string                `json:"trim"`
		Languages string                `json:"languages"`
		Whitelist string                `json:"whitelist"`
		PSM       gosseract.PageSegMode `json:"psm"`
	})

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}
	defer func(t *os.File) {
		err := tempfile.Close()
		if err != nil {
			log.Println("cannot close temp file, reason: %s", err.Error())
		}
		err = os.Remove(tempfile.Name())
		if err != nil {
			log.Println("cannot remove temp file, reason: %s", err.Error())
		}

	}(tempfile)

	if len(body.Base64) == 0 {
		render.JSON(http.StatusBadRequest, fmt.Errorf("base64 string required"))
		return
	}
	body.Base64 = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(body.Base64, "")
	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}
	tempfile.Write(b)

	client := gosseract.NewClient()
	if body.PSM > 0 {
		client.SetPageSegMode(body.PSM)
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Println("cannot close client, reason: %s", err.Error())
		}
	}()

	client.Languages = []string{"eng"}
	if body.Languages != "" {
		client.Languages = strings.Split(body.Languages, ",")
	}
	client.SetImage(tempfile.Name())
	if body.Whitelist != "" {
		client.SetWhitelist(body.Whitelist)
	}

	text, err := client.Text()
	if err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}

	render.JSON(http.StatusOK, map[string]interface{}{
		"result":  strings.Trim(text, body.Trim),
		"version": version,
	})
}
