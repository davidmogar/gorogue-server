package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorogue-server/puppet"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type Server struct{
	puppet *puppet.Puppet
}

func (server *Server) Start() error {
	registerHandlers(server)

	var err error; server.puppet, err = puppet.New()
	if err != nil {
		return err
	}

	return http.ListenAndServeTLS(fmt.Sprintf(":%d", server.puppet.MasterPort), server.puppet.HostCert, server.puppet.HostPrivKey, nil)
}

func registerHandlers(server *Server) {
	router := mux.NewRouter()

	router.HandleFunc("/puppet-ca/v1/certificate/{certname}", server.certificateHandler)
	router.HandleFunc("/puppet-ca/v1/certificate_request/{certname}", server.certificateRequestHandler)
	router.HandleFunc("/puppet-ca/v1/certificate_revocation_list/ca", server.crlHandler)
	router.HandleFunc("/puppet/v3/catalog/{certname}", catalogHandler)
	router.HandleFunc("/puppet/v3/fileContent", fileContentHandler)
	router.HandleFunc("/puppet/v3/file_metadata/{path:.*}", fileMetadataHandler)
	router.HandleFunc("/puppet/v3/file_metadatas/{path:.*}", fileMetadatasHandler)
	router.HandleFunc("/puppet/v3/node/{certname}", nodeHandler)
	router.HandleFunc("/puppet/v3/report/{certname}", reportHandler)
	router.HandleFunc("/{.*}/status/test", statusHandler)
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Use(headersMiddleware)

	http.Handle("/", router)
}

func notFound(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("NOT FOUND")
	fmt.Println(request.URL)
}

func headersMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL)
		handler.ServeHTTP(writer, request)
		writer.Header().Set("X-Puppet-Version", "rogue")
	})
}

func (server *Server) certificateHandler(writer http.ResponseWriter, request *http.Request) {
	certname := mux.Vars(request)["certname"]

	var path string
	if certname == "ca" {
		path = server.puppet.CaCert
	} else {
		path = server.puppet.SignedDir + "/" + certname
	}
	dat, _ := ioutil.ReadFile(path)
	writer.Write(dat)
}

func (server *Server) certificateRequestHandler(writer http.ResponseWriter, request *http.Request) {
	certname := mux.Vars(request)["certname"]

	file, err := os.Create(server.puppet.RequestDir + "/" + certname)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
	defer file.Close()
	io.Copy(file, request.Body)
	server.puppet.GenerateCertificate(certname)
}

func (server *Server) crlHandler(writer http.ResponseWriter, request *http.Request) {
	dat, _ := ioutil.ReadFile(server.puppet.CaCrl)
	writer.Write(dat)
}

func catalogHandler(writer http.ResponseWriter, request *http.Request) {

}

func fileContentHandler(writer http.ResponseWriter, request *http.Request) {
	//certname := mux.Vars(request)["certname"]
}

func fileMetadataHandler(writer http.ResponseWriter, request *http.Request) {
	path := mux.Vars(request)["path"]

	nodeDefinition := struct {
		Message   string    `json:"message"`
		IssueKind string    `json:"issue_kind"`
	}{
		"Not Found: Could not find file_metadata " + path,
		"RESOURCE_NOT_FOUND",
	}
	writer.Header().Set("X-Puppet-Version", "6.3.0")
	writer.Header().Add("Content-Type", "application/json;charset=utf-8")
	writer.WriteHeader(http.StatusNotFound)
	json.NewEncoder(writer).Encode(nodeDefinition)

	//metadata, err := puppet.GetFileMetadata(path)
	//if err != nil {
	//	writer.WriteHeader(http.StatusNotFound)
	//	writer.Write([]byte("{\"message\":\"Not Found: Could not find file_metadata " + path + "\",issue_kind\":\"RESOURCE_NOT_FOUND\"}"))
	//} else {
	//	json.NewEncoder(writer).Encode(metadata)
	//}
}

func fileMetadatasHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte("{}"))
}

func nodeHandler(writer http.ResponseWriter, request *http.Request) {
	nodeDefinition := struct {
		Environment string    `json:"environment"`
		Name        string    `json:"name"`
		Parameters  struct {} `json:"paramenters"`
	}{
		request.FormValue("environment"),
		mux.Vars(request)["certname"],
		struct{}{},
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(nodeDefinition)
}

func reportHandler(writer http.ResponseWriter, request *http.Request) {

}

func statusHandler(writer http.ResponseWriter, request *http.Request) {

}

