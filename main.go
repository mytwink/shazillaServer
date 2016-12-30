package main

import (
	//"database/sql"
	_ "github.com/lib/pq"
	//"io"
	"log"
	"text/template"
	"os"
	"net/http"
	//"encoding/json"
	//"strconv"
	"fmt"

	//"shazilla/mfcc"
	//"shazilla/wav"
)

var (
	dbconnect  = os.Getenv("OPENSHIFT_POSTGRESQL_DB_URL")
)

func main() {
	http.HandleFunc("/", homePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	fmt.Printf("listening on %s...", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	//wavPath := os.Getenv("OPENSHIFT_DATA_DIR")

	if r.Method == "POST" {
		/*
		r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            log.Println(err)
            return
        }
        defer file.Close()

        db, err := sql.Open("postgres", dbconnect)
		if err != nil {
			log.Println("database connection error", err)
		}
		var latestId int64
		err = db.QueryRow("SELECT id FROM samples ORDER BY id DESC LIMIT 1").Scan(&latestId)
		if err != nil {
			log.Println("Latest Id was not selected", err)
		}

        fmt.Fprintf(w, "%v", handler.Header)

        filename := strconv.FormatInt(latestId+1, 10)+"_thelaria_"+handler.Filename
        f, err := os.OpenFile(wavPath+filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            log.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
        myWavParse, err := wav.NewWavParse(file)
        if err != nil && err != io.EOF {
        	log.Println(err)
        }
        log.Println("Data parsed")
        myMfcc := mfcc.NewMfcc(myWavParse.Wav.Subchunk2.Data, myWavParse.Wav.Subchunk1.Samplerate)

        bytes, err := json.Marshal(myMfcc.GetVector())
        if err != nil {
        	log.Println("error marshalling json", err)
        }

        
		insertRes, err := db.Exec("INSERT INTO samples (type, name, vector, file_path) VALUES($1, $2, $3, $4)", "WAV", filename, string(bytes), wavPath+filename)
		if err != nil {
			log.Println("sample was not inserted", err)
		}

		insertRes.LastInsertId()

		rows, err := db.Query("SELECT id,name,file_path,vector FROM samples")
		if err != nil {
			log.Println("samples were not selected", err)
		}

		db.Close()
		samples := []mfcc.MySample{}

		for rows.Next() {
			var m mfcc.MySample
			var ror string
			rows.Scan(&m.Id, &m.Name, &m.Path, &ror)

			err = json.Unmarshal([]byte(ror), &m.Vector)
			if err != nil {
			    log.Println("error unmarshalling json", err)
			}
			samples = append(samples, m)
		}

		resSample := mfcc.MySimpleSample{}
		var minDist float64
		minDist = 1000000000.0

		var c chan mfcc.MySimpleSample = make(chan mfcc.MySimpleSample)

		for _, sample := range samples {
			if sample.Name != filename {
				vector := sample.Vector
				simple := mfcc.MySimpleSample {
					Id: sample.Id,
					Name: sample.Name,
					Path: sample.Path,
				}
				go myMfcc.Chisqr(vector, simple, c)
			}
		}

		for i:=0; i<len(samples)-1; i++ {
			res := <-c
			if res.Dist<minDist && res.Name != filename {
				log.Println(res.Id)
				minDist = res.Dist
				resSample = res
			}
		}


		data := struct {
			FileName string
			Result mfcc.MySimpleSample
		}{
			filename,
			resSample,
		}

		t, err := template.ParseFiles("templates/result.tpl")
		if err != nil {
			log.Println("template error", err)
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println("template print error", err)
		}

		db.Close()*/
	} else {
		data := struct{}{}

		t, err := template.ParseFiles("templates/index.tpl")
		if err != nil {
			log.Println("template error", err)
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println("template print error", err)
		}
	}
}