package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/serz999/notesServer/pkg/dto"
)

func main() {
    host := *flag.String("host", "localhost", "connected server host")
    port := *flag.String("port", "8000", "connected server port") 
    schema := "http" 
    notesEndpoint := "notes"
    url := fmt.Sprintf("%s://%s:%s/%s", schema, host, port, notesEndpoint)
    for true {
        Help()
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
        case "exit":
            os.Exit(0)
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
    r := bufio.NewReader(os.Stdin)
    note.Note, _  = r.ReadString('\n')
    
    jsonBytes, err := json.Marshal(note)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    
    reader := bytes.NewReader(jsonBytes)
    res, err := http.Post(noteUrl + "/", "application/json", reader)
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

    if (res.StatusCode != 200) {
        fmt.Printf("%s\n", body)
        return
    }

    var addedNote dto.Note
    json.Unmarshal(body, &addedNote)
    fmt.Printf("Added. Id of new note is '%s'\n", addedNote.Id)
}

func Del(noteUrl string) {
    fmt.Printf("Enter note id: ") 
    var id dto.Id
    fmt.Scanf("%s", &id)

    fmt.Printf("Are you shure what to delete this note?(f): ")
    var answer string
    fmt.Scanf("%s", &answer)
    if !Yes(answer) {
        return 
    }

    client := http.Client{}
    req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", noteUrl, id), nil) 
    if err != nil {
        fmt.Println(err.Error()) 
        return
    }

    res, err := client.Do(req)
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
    
    if res.StatusCode == 404 {
        fmt.Println("Note with such id does not exist.") 
        return
    } else if (res.StatusCode != 200) {
        fmt.Printf("%s\n", body)
        return
    } 

    fmt.Println("Deleted.")
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

    if res.StatusCode == 404 {
        fmt.Println("Note with such id does not exist.") 
        return
    } else if (res.StatusCode != 200) {
        fmt.Printf("%s\n", body)
        return
    }
    
    var note dto.Note 
    deserializeerr := json.Unmarshal(body, &note)
    if deserializeerr != nil {
        fmt.Println(deserializeerr.Error())
        return
    } 

    NoteView(note)
}

func NoteView(note dto.Note) {
    border := ""
    for i := 0; i < len(note.Note); i++ {
        border += "-" 
    }

    fmt.Println(border)
    fmt.Printf("%s", note.Note)
    fmt.Println(border)

    fmt.Printf("Text by: %s %s\n", note.AuthorFirstName, note.AuthorLastName)
}

func Help() {
    fmt.Println("PROMPT")
    fmt.Println("   help    - print help")
    fmt.Println("   add     - add note")
    fmt.Println("   del     - del note")
    fmt.Println("   get     - get note")
    fmt.Println("   exit    - stop the program execution")
}
