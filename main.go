package main

import (
	"io"
    "os"
    "log"
	"io/ioutil"
	"net/http"
	"errors"
	//"fmt"
	"strconv"
	"strings"
	"encoding/json"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
)

func UnmarshalCharacter(data []byte) (Character, error) {
	var r Character
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Character) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Character struct {
	Info    Info     `json:"info"`   
	Results []Result `json:"results"`
}

type Info struct {
	Count int64  `json:"count"`
	Pages int64  `json:"pages"`
	Next  string `json:"next"` 
	Prev  string `json:"prev"` 
}

type Result struct {
	ID       int   `json:"id"`      
	Name     string   `json:"name"`    
	Status   Status   `json:"status"`  
	Species  string   `json:"species"` 
	Type     string   `json:"type"`    
	Gender   Gender   `json:"gender"`  
	Origin   Location `json:"origin"`  
	Location Location `json:"location"`
	Image    string   `json:"image"`   
	Episode  []string `json:"episode"` 
	URL      string   `json:"url"`     
	Created  string   `json:"created"` 
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"` 
}

type Gender string
const (
	Female Gender = "Female"
	GenderUnknown Gender = "unknown"
	Male Gender = "Male"
)

type Status string
const (
	Alive Status = "Alive"
	Dead Status = "Dead"
	StatusUnknown Status = "unknown"
)

var chars []Result
var forms []*widget.Form

func main() {
	myApp := app.New()
	logo, err := fyne.LoadResourceFromPath("logo4.jpg")
	if err != nil{
		log.Println("error during loading logo")
	}
	myApp.SetIcon(logo)

	pages := []string{"1", "2", "3", "4", "5", "6"}
	for _, num := range pages{

		resp, err := http.Get("https://rickandmortyapi.com/api/character/?page="+num)
		if err != nil{
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil{
			log.Fatalln(err)
		}

		char, err := UnmarshalCharacter(body)
		if err != nil{
			log.Println("error during unmarshall")
		}

		for _, be := range char.Results{
			chars = append(chars, be)
		}
	}

	form := widget.NewForm()
	form.Append("ID", widget.NewLabel("Character's Name"))
	
	for _, be := range chars{
		id := be.ID
		form.AppendItem(widget.NewFormItem(strconv.Itoa(be.ID), widget.NewButton(be.Name, func(){ShowCharacter(id, myApp)})))
	}
	forms = append(forms, form)

	myWindow := myApp.NewWindow("Rick&Morty")

	myWindow.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("Menu",
		fyne.NewMenuItem("Dark Theme", func() { myApp.Settings().SetTheme(theme.DarkTheme()) }),
		fyne.NewMenuItem("Light Theme", func() {myApp.Settings().SetTheme(theme.LightTheme()) }),
	)))

	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Character's Name: ")

	someSpace := widget.NewLabel("")

	searchButton := widget.NewButton("Search", func() {
		form1, err := SearchCharacter(input.Text, myApp)
		if err != nil{
			dialog.ShowError(err, myWindow)
		}else{
			forms = append(forms, form1)
			SearchResults(myApp)
		}
	})

	content := widget.NewVBox(input, searchButton, someSpace, GetLastForm())
	

	
	myWindow.SetContent(widget.NewScrollContainer(content))
	myWindow.SetMaster()
	myWindow.Resize(fyne.NewSize(form.MinSize().Width+50, 480))
	//myWindow.SetFixedSize(true)
	myWindow.ShowAndRun()
}


func GetLastForm() *widget.Form{
	return forms[len(forms)-1]
}

func SearchResults(myApp fyne.App){
	resWindow := myApp.NewWindow("Search Results")

	form := GetLastForm()
	resWindow.SetContent(fyne.NewContainerWithLayout(layout.NewFormLayout(), form))
	resWindow.Resize(fyne.NewSize(form.MinSize().Width+200, form.MinSize().Height+100))
	resWindow.SetFixedSize(true)
	resWindow.Show()
}

func SearchCharacter(input string, myApp fyne.App) (*widget.Form, error){

	form := widget.NewForm()
	form.Append("ID", widget.NewLabel("Character's Name"))
	counter := 0
	for _, be := range chars{
		if strings.Contains(strings.ToLower(be.Name), strings.ToLower(input)){
			id := be.ID
			counter++
			form.AppendItem(widget.NewFormItem(strconv.Itoa(be.ID), widget.NewButton(be.Name, func(){ShowCharacter(id, myApp)})))
		}
	}
	//log.Println(counter)
	err := errors.New("Character not found")
	if counter != 0{
		err = nil
		
	}
	if counter == 1{

	}
	//log.Println("search's done")
	return form, err
}



func ShowCharacter(id int, myApp fyne.App){
	//log.Println(id)
	var char Result

	for _, be := range chars{
		if id == be.ID{
			char = be
		}
	}

	minimain(strconv.Itoa(id))

	charWindow := myApp.NewWindow(char.Name)
	//charWindow.SetIcon("")
	img := canvas.NewImageFromFile("CharacterImages/"+strconv.Itoa(id)+".jpg")
	eps := ""
	for _, ep := range char.Episode{
			eps += ep
	}

	form := widget.NewForm()
	form.Append("ID: ", widget.NewLabel(strconv.Itoa(id)))
	form.Append("Name: ", widget.NewLabel(char.Name))
	form.Append("Status: ", widget.NewLabel(string(char.Status)))
	form.Append("Species: ", widget.NewLabel(char.Species))
	form.Append("Type: ", widget.NewLabel(char.Type))
	form.Append("Gender: ", widget.NewLabel(string(char.Gender)))
	form.Append("Origin: ", widget.NewLabel(char.Origin.Name))
	form.Append("Location: ", widget.NewLabel(char.Location.Name))

	//log.Println("Here is the min Size: ")
	//log.Println(form.MinSize())

	charWindow.SetContent(fyne.NewContainerWithLayout(layout.NewFormLayout(),
		form,
		img))

	charWindow.Resize(fyne.NewSize(form.MinSize().Width+300, 380))
	charWindow.SetFixedSize(true)
	charWindow.Show()
}


var (
    fileName    string
    fullUrlFile string
)


func minimain(id string){


	fullUrlFile = "https://rickandmortyapi.com/api/character/avatar/"+id+".jpeg"

    fileName = "CharacterImages/"+id+".jpg"

    // Create blank file
    file := createFile()

    // Put content on file
    putFile(file, httpClient())
}

func putFile(file *os.File, client *http.Client) {
    resp, err := client.Get(fullUrlFile)

    checkError(err)

    defer resp.Body.Close()

    io.Copy(file, resp.Body)

    defer file.Close()

    checkError(err)

    
}

func httpClient() *http.Client {
    client := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
        },
    }

    return &client
}

func createFile() *os.File {
    file, err := os.Create(fileName)

    checkError(err)
    return file
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}