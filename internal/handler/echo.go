package handler

import (
	"net/http"
	"io"
	"log"
)
func EchoHandler(writer http.ResponseWriter, req *http.Request) {

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
