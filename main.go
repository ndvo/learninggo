package main

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

// Folders
var folders map[string]string = map[string]string{
	"app":      "app",
	"api":      "api",
	"css":      "css",
	"users":    "users",
	"content":  "content",
	"delta":    "delta",
	"template": "api/templates"}

// init creates the folder structure and databases necessary for dom
func init() {
	println("Criando os diretórios necessários:")
	for _, f := range folders {
		println("\tCriando " + f)
		os.MkdirAll(f, 0775)
	}
	testuser := &Usuario{
		saver:    saver{addr: "nelson"},
		username: "nelson",
		password: []byte("teste de senha"),
		email:    "teste+teste@ocastudios.com",
	}
	println("Teste de usuario")
	println(testuser.addr)
	println(testuser.username)
	println("now let's try to save")
	testuser.save()
}

// Creates a saver struct that has the address to witch it should be saved
type saver struct {
	addr string
}

// save gzips and saves a struct s to a file addr
// It does not create the directories and assumes the permissions are right
func (s *saver) save() error {
	println("... tentando salvar " + s.addr)
	f, err := os.Create(s.addr)
	if err != nil {
		return err
	}
	defer f.Close()
	fz := gzip.NewWriter(f)
	defer fz.Close()
	enc := gob.NewEncoder(fz)
	err = enc.Encode(s.addr)
	if err != nil {
		return err
	}
	return nil
}

type session struct {
	Usuario
	sid string
}

// createID creates a session ID from system random and returns a sessionIDCreate a new session id
func (s *session) createID() {
	//32 bytes == 256 bits
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		println(err)
	}
	s.sid = base64.URLEncoding.EncodeToString(b)
}

type Imobiliaria struct {
	*saver
	nome string
}

type Vizinhança struct {
	*saver
	caracteristicas string
}

type Condominio struct {
	*saver
}

type Anuncio struct {
	*saver
}

type Imovel struct {
	*saver
	tipo        string
	endereço    Endereco
	preço       string
	dir         string
	imobiliaria Imobiliaria
	vizinhanca  string
}

type Endereco struct {
	cep, estado, cidade, bairro, rua, latitude, longitude, descritivo string
}

type Usuario struct {
	saver
	username string
	password []byte
	email    string
	addr     string
}

type Agente struct {
	usuario      Usuario
	representado Usuario
	imobiliaria  Imobiliaria
}

type Page struct {
	Title string
	Body  template.HTML
	Delta string
	User  string
}

func (p *Page) save() {
	go ioutil.WriteFile("content/"+p.Title, []byte(p.Body), 0666)
	go ioutil.WriteFile("delta/"+p.Title+".delta", []byte(p.Delta), 0666)
}

type PageInfo struct {
	FileInfo os.FileInfo
	Name     string
	Size     int64
	ModTime  time.Time
}

func (p *PageInfo) FileName() string {
	nome := p.FileInfo.Name()
	return basename(nome)
}

func basename(fileName string) string {
	re := regexp.MustCompile(`^([\p{L}_\d]+)\.\w+$`)
	matched := re.FindStringSubmatch(fileName)
	if matched == nil {
		return "Não encontrado"
	} else {
		return matched[1]
	}
}

func editPage(w http.ResponseWriter, r *http.Request) {
	title := strings.Split(r.URL.Path, "/")[2]
	println("Editing Page " + title)
	//cookie, err := r.Cookie("user")
	//var user string
	//if err != nil {
	//	user = cookie.Value
	//} else {
	//	user = ""
	//}
	//println("Afer knowing cookie")
	page, err := loadPage(title, "")
	if err != nil {
		page = &Page{Title: title, Body: "", Delta: "", User: "TODO"}
	}
	t, _ := template.ParseFiles("api/templates/edit_page.html")
	t.Execute(w, page)
}

func loadPage(title string, user string) (*Page, error) {
	var delta []byte
	delta, err := ioutil.ReadFile("delta/" + title + ".delta")
	var delta_string string
	if err != nil {
		delta_string = ""
	} else {
		delta_string = string(delta)
	}
	body, err := ioutil.ReadFile("content/" + title)
	var body_string template.HTML
	if err != nil {
		body_string = ""
	} else {
		body_string = template.HTML(body)
	}
	return &Page{
		Title: title,
		Body:  body_string,
		Delta: delta_string,
		User:  user,
	}, nil
}

func viewPage(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	var user string
	var frontpage bool
	if title == "" {
		frontpage = true
	}
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("user")
		if err == nil {
			user = cookie.Value
		} else {
			user = ""
		}
		if frontpage == true {
			frontPage(w, r)
		} else {
			page, err := loadPage(title, user)
			if err != nil {
				http.NotFound(w, r)
				fmt.Print(err)
			} else {
				t, _ := template.ParseFiles("api/templates/view_page.html")
				t.Execute(w, page)
			}
		}
	}
}

func listContent() []PageInfo {
	pages := []PageInfo{}
	arquivos, _ := ioutil.ReadDir("content")
	for _, arq := range arquivos {
		p := PageInfo{
			arq,
			arq.Name(),
			arq.Size(),
			arq.ModTime(),
		}
		println(arq.Name())
		pages = append(pages, p)
	}
	return pages
}

func frontPage(w http.ResponseWriter, r *http.Request) {
	arquivos, err := ioutil.ReadDir("content")
	if err == nil {
		var pages []PageInfo
		for _, arq := range arquivos {
			p := PageInfo{
				arq,
				arq.Name(),
				arq.Size(),
				arq.ModTime(),
			}
			pages = append(pages, p)
		}
		frontbody, _ := ioutil.ReadFile("api/templates/front_page.html")
		print(frontbody)
		//t,_:=template.ParseFiles("api/templates/front_page.html")
		fmt.Fprint(w, string(frontbody))
		//t.Execute(w,pages)
	} else {
		print(err.Error() + "\n")
	}
}

func savePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := strings.Split(r.URL.Path, "/")[2]
		body := template.HTML(r.FormValue("body"))
		p := &Page{Title: title, Body: body, Delta: r.FormValue("delta")}
		p.save()
		fmt.Fprint(w, "documento salvo")
		//http.Redirect(w,r,"/"+title, http.StatusFound)
	}
}

type Profile struct {
	Name string
	Type string
}

func api(w http.ResponseWriter, r *http.Request) {
	print("\nAPI REQUEST: list content \n")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(listContent())
}

func loginForm(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		login_page, _ := ioutil.ReadFile("api/templates/login.html")
		fmt.Fprint(w, string(login_page))
	case "POST":
		r.ParseForm()
		userName := r.Form["nome"][0]
		println(r.Form["senha"][0])
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(r.Form["senha"][0]), 14)
		novousuario := r.Form["novousuario"][0]
		email := r.FormValue("email")
		println(novousuario)
		println(email)
		println(hashPassword)
		var u Usuario
		if novousuario == "on" {
			u = Usuario{
				saver:    saver{userName},
				username: userName,
				password: hashPassword,
				email:    r.FormValue("email"),
				addr:     userName,
			}
			println("Novo usuário " + u.username + " (" + u.email + ") criado.")
			println("Novo usuário ... tentando salvar")
			u.save()
		}
		newsession := session{
			Usuario: u,
		}
		newsession.createID()
		cookie := http.Cookie{
			Name:     "user",
			Value:    newsession.sid,
			HttpOnly: true,
			Expires:  time.Now().Add(120 * time.Second),
			MaxAge:   0,
			Path:     "/",
		}
		println("Cookie criado")
		http.SetCookie(w, &cookie)
		println("Cookie atribuído")
		http.Redirect(w, r, "/", http.StatusFound)
		println("Redirecionamento")
	}
}

func main() {
	http.HandleFunc("/", viewPage)
	http.HandleFunc("/knock-knock/", loginForm)
	http.HandleFunc("/api/", api)
	http.HandleFunc("/edit/", editPage)
	http.HandleFunc("/save/", savePage)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("app"))))
	http.Handle("/vc_components/", http.StripPrefix("/vc_components/", http.FileServer(http.Dir("app/vc_components"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
	println("Server is up and running.")
}
