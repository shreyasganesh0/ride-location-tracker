package main

import (
	"log"
	"io"
	"net/http"
)

func echo_handler(writer http.ResponseWriter, req *http.Request) {

	log.Println("In the echo_handler");
	buf := make([]byte, 1024);
	resp_buf := make([]byte, 0);
	for {
	
		n, err := req.Body.Read(buf);
		if n > 0 {
			buf = buf[:n];
			log.Printf("adding bytes %q", buf);
			resp_buf = append(resp_buf, buf...);
		}

		if (err == io.EOF) {

			log.Println("Finished reading.");
			break;
		}
		if err != nil {
			log.Println("Found Error while reading");
			log.Fatal(err.Error());
		}

	}



	log.Printf("writing bytes %q", resp_buf);
	_, err := writer.Write(resp_buf)
	if err != nil {
		log.Fatal(err.Error());
	}
}

func main() {
	log.Println("Starting echo server...");

	def_handler := func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("This is the default response."));
		if err != nil {
			log.Fatal(err.Error());
		}
	}

	http.HandleFunc("/" , def_handler);
	http.HandleFunc("/echo" , echo_handler);
	log.Fatal(http.ListenAndServe(":8080", nil));
}


