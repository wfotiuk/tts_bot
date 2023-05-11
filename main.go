package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

type SpeechBody struct {
	Sentence string
}

func main() {
	fmt.Println("program started")
	http.HandleFunc("/", handleSpeech)
	err := http.ListenAndServe(":443", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("server started")
}

func handleSpeech(w http.ResponseWriter, req *http.Request) {
	sentence := "default sentence"
	if req.Method == http.MethodPost {
		speechBody := SpeechBody{}
		err := json.NewDecoder(req.Body).Decode(&speechBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sentence = speechBody.Sentence
	}

	speech := tts.Speech{
		Folder:   "audio",
		Language: voices.English,
		Handler:  &handlers.Native{},
	}

	fmt.Printf("starting to speak sentence: %s \n\n", sentence)
	err := speech.Speak(sentence)
	if err != nil {
		fmt.Printf("error speaking %s", err)
	}

	// remove audio directory
	err = os.RemoveAll("audio")
	if err != nil {
		fmt.Printf("error removing directory %s", err)
	}

	fmt.Print("finished speaking\n\n")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("speaking shit complete.\n"))
}
