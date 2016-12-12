// http
package http

import (
	"stats/char"
	"stats/word"
	"fmt"
	"net/http"
	"strconv"
	"html/template"
)

type OutHttpServer struct {
	CharStats *char.CharStats
	WordStats *word.WordStats
	Route string
	Port int
}

func New(route string, port int, charStats *char.CharStats, wordStats *word.WordStats) (server *OutHttpServer, err error) {
	server = new(OutHttpServer)
	server.Route = route
	server.Port = port
	server.CharStats = charStats
	server.WordStats = wordStats
	
	return
}

func (s *OutHttpServer) Listen() {
	fmt.Printf("HTTP interface listening in: :%d%s\n", s.Port, s.Route)
	http.HandleFunc(s.Route, s.handleHTTPRequests)
	http.ListenAndServe(":" + strconv.Itoa(s.Port), nil)
}

type outputData struct {
	N int
	Count uint64
	Words []string
	Chars []string
}
const outputTemplateStr = `{
  "count": {{.Count}},
  "top_{{.N}}_words": [{{range $i, $e := .Words}}{{if $i}},{{end}}"{{$e}}"{{end}}],
  "top_{{.N}}_letters": [{{range $i, $e := .Chars}}{{if $i}},{{end}}"{{$e}}"{{end}}]
}`

func (s *OutHttpServer)handleHTTPRequests(w http.ResponseWriter, r *http.Request){
	//Geting the param into the URL or the Form
	n, err := strconv.ParseInt(r.URL.Query().Get("N"), 10, 32)
	if err != nil {
		n, err = strconv.ParseInt(r.FormValue("N"), 10, 32) 
		if err != nil {
			n = 5
		}
	}
		
	chars := s.CharStats.GetFirstN(uint64(n))
	newChars := make([]string, len(chars))
	for i := 0; i < len(chars); i++ {
		newChars[i] = string([]uint8{ chars[i] })
	}
	data := outputData{
		N: int(n),
		Count: s.WordStats.Count,
		Words: s.WordStats.GetFirstN(uint64(n)),
		Chars: newChars,
	}
	
	header := w.Header()
	header["Content-Type"] = []string{ "application/json" }
	tmp, _ := template.New("JSON").Parse(outputTemplateStr)
	tmp.Execute(w, data)
}