package serve

import (
	"html/template"
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

func RunServer() {
	http.HandleFunc("/", llamarManejador(manejadorRaiz)) //Nuevo manejador
	http.HandleFunc("/view/", llamarManejador(manejadorMostrar))
	http.HandleFunc("/save/", llamarManejador(manejadorGuardar))
	http.HandleFunc("/edit/", llamarManejador(manejadorEditar))

	http.ListenAndServe(":5444", nil)
}
