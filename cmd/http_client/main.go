package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ananikitina/http/models"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
)

const (
	baseUrl       = "http://localhost:8080"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

func createNote() (models.Note, error) {
	note := &models.NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}
	data, err := json.Marshal(note)
	if err != nil {
		return models.Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return models.Note{}, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return models.Note{}, err
	}

	var createdNote models.Note

	if err = json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return models.Note{}, err
	}
	return createdNote, nil
}

func getNote(id int64) (models.Note, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl+getPostfix, id))
	if err != nil {
		log.Fatal("Failed to get note", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return models.Note{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return models.Note{}, err
	}
	var note models.Note
	if err = json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal("Failed to create note", err)
	}
	log.Printf(color.RedString("Note created:\n"), color.GreenString("%+v", note))

	note, err = getNote(note.ID)
	if err != nil {
		log.Fatal("Failed to get note", err)
		return
	}

	log.Printf(color.RedString("Note info got:\n"), color.GreenString("%+v", note))
}
