package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"github.com/serz999/notesServer/pkg/dto"
)

func main() {
    host := *flag.String("host", "localhost", "connected server host")
    port := *flag.String("port", "8000", "connected server port") 
    schema := "http" 
    notesEndpoint := "notes"
    url := fmt.Sprint("%s://%s:%s/%s", schema, host, port, notesEndpoint) 
    for true {
        fmt.Printf("note> ")
        var cmd string
        fmt.Scanf("%s", &cmd)
        switch cmd {
        case "add":
            Add(url)
        case "del":
            Del(url)
        case "get":
            Get(url)
        case "help":
            Help()
        default:
            fmt.Println("Invalid command, try 'help' prompt")
        } 
    }
}

func Add(noteUrl string) {
    note := dto.Note{} 

    fmt.Printf("Enter the firstname: ")
    fmt.Scanf("%s", &note.AuthorFirstName) 
    fmt.Printf("Enter the lastname: ")
    fmt.Scanf("%s", &note.AuthorLastName) 
    fmt.Printf("Your note:")
    fmt.Scanf("%s", &note.Note)
    
    jsonBytes, err := json.Marshal(note)
    if err != nil {
        fmt.Println(err.Error())
    }

    reader := bytes.NewReader(jsonBytes)
    res, err := http.Post(noteUrl, "application/json", reader)
    if err != nil {
        fmt.Println(err.Error())
    } 
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err.Error())
    }

    var addedNote dto.Note
    json.Unmarshal(body, &addedNote)
    fmt.Printf("Added. Id of new note is '%v'", addedNote.Id)
}

func Del(noteUrl string) {
    fmt.Printf("Enter note id: ") 
    var id dto.Id
    fmt.Scanf("%s", id)

    fmt.Printf("Are you shure what to delete this note?(f): ")
    var answer string
    fmt.Scanf("%s", &answer)
    if !Yes(answer) {
        return      
    }

    res, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", noteUrl, id), nil) 
    if err != nil {
        fmt.Printf(err.Error()) 
        return
    }
    defer res.Body.Close()

    fmt.Print("Deleted.")
}

func Yes(answer string) bool {
    return answer == "y" || answer == "Y" || answer == "yes"
}

func Get(noteUrl string) {
    fmt.Printf("id: ")
    var id dto.Id 
    fmt.Scanf("%s", &id)

    res, err := http.Get(fmt.Sprintf("%s/%s", noteUrl, id))
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    defer res.Body.Close()
    
    body, err := io.ReadAll(res.Body)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    var note dto.Note 
    deserializeerr := json.Unmarshal(body, note)
    if deserializeerr != nil {
        fmt.Println(deserializeerr.Error())
        return
    } 

    fmt.Printf("%+v", note)
}

func Help() {
    fmt.Println("SCRIPTS")
    fmt.Println("   help - print help")
    fmt.Println("   add - add note")
    fmt.Println("   del - del note")
    fmt.Println("   get - get note")
}
