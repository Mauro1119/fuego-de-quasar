package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type responseDecoded struct {
	Pos     position `json:"position"`
	Message string   `json:"message"`
}

type position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type satellites struct {
	Satellite []satellite `json:"satellites"`
}

type satellite struct {
	Name     string   `json:"name"`
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

type onlyOneSatellite struct {
	Distance float32  `json:"distance"`
	Message  []string `json:"message"`
}

var sat = satellites{}

func getSatellites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sat)
}

func createSatellites(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid satellite info")
		return
	}
	json.Unmarshal(reqBody, &sat)

	resp, errSat := Procesar()
	w.Header().Set("Content-Type", "application/json")
	if resp.Message == "" || errSat != 0 {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func addSatellite(w http.ResponseWriter, r *http.Request) {
	var newSatellite satellite
	var newSat onlyOneSatellite
	var nuevo bool = true
	vars := mux.Vars(r)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Invalid satellite info")
		return
	}
	json.Unmarshal(reqBody, &newSat)

	newSatellite.Name = vars["satellite_name"]
	newSatellite.Distance = newSat.Distance
	newSatellite.Message = newSat.Message

	for i, s := range sat.Satellite {
		if strings.ToUpper(s.Name) == strings.ToUpper(newSatellite.Name) {
			sat.Satellite = append(sat.Satellite[:i], sat.Satellite[i+1:]...)
			newSatellite.Name = s.Name
			nuevo = false
		}
	}
	if nuevo {
		if len(sat.Satellite) == 3 {
			fmt.Fprintf(w, "Demasiados satellites")
			return
		}
	}

	sat.Satellite = append(sat.Satellite, newSatellite)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(newSatellite)
}

func getPosition(w http.ResponseWriter, r *http.Request) {
	resp, err := Procesar()
	w.Header().Set("Content-Type", "application/json")
	if resp.Message == "" || err != 0 {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func deleteAll(w http.ResponseWriter, r *http.Request) {
	sat.Satellite = sat.Satellite[:0]
	w.WriteHeader(http.StatusOK)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/topsecret", createSatellites).Methods("POST") //retorna posicion y mensaje si se puede, sino 404
	router.HandleFunc("/topsecret", getSatellites).Methods("GET")     //retorna todos los satellites
	router.HandleFunc("/topsecret", deleteAll).Methods("DELETE")
	router.HandleFunc("/topsecret_split/{satellite_name}", addSatellite).Methods("POST") //crea un satellite
	router.HandleFunc("/topsecret_split", getPosition).Methods("GET")                    //retorna posicion y mensaje si se puede

	log.Fatal(http.ListenAndServe(":8080", router))

}
