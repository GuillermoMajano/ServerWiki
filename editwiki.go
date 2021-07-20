package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type Pagina struct {
	Titulo string
	Cuerpo []byte
}

var plantillas = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html", "tmpl/front.html"))
var regex_ruta = regexp.MustCompile("^(/|(/(edit|save|view)/([a-zA-Z0-9]+)))$")
var pagina_principal = "Principal"

func main() {
	http.HandleFunc("/", llamarManejador(manejadorRaiz)) //Nuevo manejador
	http.HandleFunc("/view/", llamarManejador(manejadorMostrar))
	http.HandleFunc("/save/", llamarManejador(manejadorGuardar))
	http.HandleFunc("/edit/", llamarManejador(manejadorEditar))

	http.ListenAndServe(":5444", nil)
}

//Método para guardar página
func (p *Pagina) guardar() error {
	nombre := p.Titulo + ".txt"
	return ioutil.WriteFile("./view/"+nombre, p.Cuerpo, 0600)
}

//Función para cargar página
func cargarPagina(titulo string) (*Pagina, error) {
	nombre_archivo := titulo + ".txt"
	fmt.Println("El cliente ha pedido: " + nombre_archivo)
	cuerpo, err := ioutil.ReadFile("./view/" + nombre_archivo)
	if err != nil {
		return nil, err
	}
	return &Pagina{Titulo: titulo, Cuerpo: cuerpo}, nil
}

//Función para validar ruta y regresar nombre de la página solicitada
func dameTitulo(w http.ResponseWriter, r *http.Request) (string, error) {
	peticion := regex_ruta.FindStringSubmatch(r.URL.Path)
	if peticion == nil {
		http.NotFound(w, r)
		return "", errors.New("Ruta inválida")
	}
	return peticion[len(peticion)-1], nil
}

//Función para cargar las plantillas HTML
func cargarPlantilla(w http.ResponseWriter, nombre_plantilla string, pagina *Pagina) {
	plantillas.ExecuteTemplate(w, nombre_plantilla+".html", pagina)
}
func llamarManejador(manejador func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		titulo, err := dameTitulo(w, r)
		fmt.Println(titulo)
		if err != nil {
			return
		}
		manejador(w, r, titulo)
	}
}

//Manejador para mostrar página principal
func manejadorRaiz(w http.ResponseWriter, r *http.Request, titulo string) {
	p, err := cargarPagina(pagina_principal)
	if err != nil {
		http.Redirect(w, r, "edit/"+pagina_principal, http.StatusFound)
		fmt.Println("Error")
		return
	}
	cargarPlantilla(w, "front", p)
}

//Manejador para visualizar wikis
func manejadorMostrar(w http.ResponseWriter, r *http.Request, titulo string) {
	p, err := cargarPagina(titulo)
	if err != nil {
		http.Redirect(w, r, "/edit/"+titulo, http.StatusFound)
		fmt.Println("La página solicitada no existía. Llamando al editor...")
		return
	}
	cargarPlantilla(w, "view", p)
}

//Manejador para editar wikis
func manejadorEditar(w http.ResponseWriter, r *http.Request, titulo string) {
	p, err := cargarPagina(titulo)
	if err != nil {
		p = &Pagina{Titulo: titulo}
	}
	cargarPlantilla(w, "edit", p)
}

//Manejador para guardar wikis
func manejadorGuardar(w http.ResponseWriter, r *http.Request, titulo string) {
	cuerpo := r.FormValue("body")
	p := &Pagina{Titulo: titulo, Cuerpo: []byte(cuerpo)}
	fmt.Println("Guardando " + titulo + ".txt...")
	p.guardar()
	http.Redirect(w, r, "/view/"+titulo, http.StatusFound)
}
