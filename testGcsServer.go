package main

import (
	"appengine"
	// "appengine/datastore"
	"appengine/file"
	// "bytes"
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	// "strings"
	"time"

	"github.com/pborman/uuid"
	// "golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	newappengine "google.golang.org/appengine"
	// "google.golang.org/appengine/file"
	newurlfetch "google.golang.org/appengine/urlfetch"
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
)

// type Result struct {
// 	ReturnCode int `json:"returncode"`
// }

type Book struct {
	Name       string    `json:"name"`
	Author     string    `json:"author"`
	Pages      int       `json:"pages"`
	Year       int       `json:"year"`
	CreateTime time.Time `json:"createtime"`
}

const BaseUrl = "/api/0.1/"
const BookKind = "Book"
const BookRoot = "Book Root"
const BookMaxPages = 1000

var BookName = []string{"AAA", "BBB", "CCC", "DDD", "EEE", "FFF", "GGG", "HHH", "III", "JJJ"}
var BookAuthor = []string{"AuthorA", "AuthorB", "AuthorC", "AuthorD", "AuthorE", "AuthorF", "AuthorG", "AuthorH", "AuthorI", "AuthorJ"}

func init() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc(BaseUrl, rootPage)
	// http.HandleFunc(BaseUrl+"queryAll", queryAll)
	// http.HandleFunc(BaseUrl+"queryAllWithKey", queryAllWithKey)
	http.HandleFunc(BaseUrl+"storeImage", storeImage)
	// http.HandleFunc(BaseUrl+"deleteAll", deleteAll)
	// http.HandleFunc(BaseUrl+"images", images)
	log.Println("I'm out init()")
}

func rootPage(rw http.ResponseWriter, req *http.Request) {
	//
}

// func images(rw http.ResponseWriter, req *http.Request) {
// 	switch req.Method {
// 	case "GET":
// 		queryBook(rw, req)
// 	case "POST":
// 		storeImage(rw, req)
// 	case "DELETE":
// 		deleteBook(rw, req)
// 	default:
// 		queryAll(rw, req)
// 	}
// }

func storeImage(rw http.ResponseWriter, req *http.Request) {
	// Result, 0: success, 1: failed
	var r int = 0
	fileName := uuid.New()+".jpg"
	// var cKey *datastore.Key = nil
	defer func() {
		// Return status. WriteHeader() must be called before call to Write
		if r == 0 {
			// Changing the header after a call to WriteHeader (or Write) has no effect.
			// rw.Header().Set("Location", req.URL.String()+"/"+cKey.Encode())
			rw.Header().Set("Location", req.URL.String()+"/"+fileName)
			rw.WriteHeader(http.StatusCreated)
		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	}()

	// Get data from body
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintln(rw, err, "in reading body")
		r = 1
		return
	}
	fmt.Fprintln(rw, "Body length:", len(b))

	cc := appengine.NewContext(req)
	bucket := ""
	if bucket, err = file.DefaultBucketName(cc); err != nil {
		fmt.Fprintln(rw, err, "in getting default GCS bucket name")
		r = 1
		return
	}
	c := newappengine.NewContext(req)
	hc := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(c, storage.ScopeReadWrite),
			// Note that the App Engine urlfetch service has a limit of 10MB uploads and
			// 32MB downloads.
			// See https://cloud.google.com/appengine/docs/go/urlfetch/#Go_Quotas_and_limits
			// for more information.
			Base: &newurlfetch.Transport{Context: c},
		},
	}
	ctx := cloud.NewContext(newappengine.AppID(c), hc)
	fmt.Fprintln(rw, "APP Engine Version:", newappengine.VersionID(c))
	fmt.Fprintln(rw, "Using bucket name:", bucket)

	// d := &demo{
	// 	c:   c,
	// 	w:   w,
	// 	ctx: ctx,
	// }

	// createFile creates a file in Google Cloud Storage.
	fmt.Fprintf(rw, "Creating file /%v/%v\n", bucket, fileName)

	wc := storage.NewWriter(ctx, bucket, fileName)
	wc.ContentType = "image/jpeg"
	// wc.Metadata = map[string]string{
	// 	"x-goog-meta-foo": "foo",
	// 	"x-goog-meta-bar": "bar",
	// }
	if _, err := wc.Write(b); err != nil {
		fmt.Fprintln(rw, "createFile: unable to write data to bucket %q, file %q: %v", bucket, fileName, err)
		r = 1
		return
	}
	if err := wc.Close(); err != nil {
		fmt.Fprintln(rw, "createFile: unable to close bucket %q, file %q: %v", bucket, fileName, err)
		r = 1
		return
	}
}

// func queryAll(rw http.ResponseWriter, req *http.Request) {
// 	// Get all entities
// 	var dst []Book
// 	r := 0
// 	c := appengine.NewContext(req)
// 	_, err := datastore.NewQuery(BookKind).Order("Pages").GetAll(c, &dst)
// 	if err != nil {
// 		log.Println(err)
// 		r = 1
// 	}

// 	// Return status. WriteHeader() must be called before call to Write
// 	if r == 0 {
// 		rw.WriteHeader(http.StatusOK)
// 	} else {
// 		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	// Return body
// 	encoder := json.NewEncoder(rw)
// 	if err = encoder.Encode(dst); err != nil {
// 		log.Println(err, "in encoding result", dst)
// 	} else {
// 		log.Printf("QueryAll() returns %d items\n", len(dst))
// 	}
// }

// func queryBook(rw http.ResponseWriter, req *http.Request) {
// 	if len(req.URL.Query()) == 0 {
// 		queryAllWithKey(rw, req)
// 	} else {
// 		searchBook(rw, req)
// 	}
// }

// func queryAllWithKey(rw http.ResponseWriter, req *http.Request) {
// 	// Get all entities
// 	var dst []Book
// 	r := 0
// 	c := appengine.NewContext(req)
// 	k, err := datastore.NewQuery(BookKind).Order("Pages").GetAll(c, &dst)
// 	if err != nil {
// 		log.Println(err)
// 		r = 1
// 	}

// 	// Map keys and books
// 	var m map[string]*Book
// 	m = make(map[string]*Book)
// 	for i := range k {
// 		m[k[i].Encode()] = &dst[i]
// 	}

// 	// Return status. WriteHeader() must be called before call to Write
// 	if r == 0 {
// 		rw.WriteHeader(http.StatusOK)
// 	} else {
// 		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	// Return body
// 	encoder := json.NewEncoder(rw)
// 	if err = encoder.Encode(m); err != nil {
// 		log.Println(err, "in encoding result", m)
// 	} else {
// 		log.Printf("QueryAll() returns %d items\n", len(m))
// 	}
// }

// func searchBook(rw http.ResponseWriter, req *http.Request) {
// 	// Get all entities
// 	var dst []Book
// 	r := 0
// 	q := req.URL.Query()
// 	f := datastore.NewQuery(BookKind)
// 	for key := range q {
// 		f = f.Filter(key+"=", q.Get(key))
// 	}
// 	c := appengine.NewContext(req)
// 	k, err := f.GetAll(c, &dst)
// 	if err != nil {
// 		log.Println(err)
// 		r = 1
// 	}

// 	// Map keys and books
// 	var m map[string]*Book
// 	m = make(map[string]*Book)
// 	for i := range k {
// 		m[k[i].Encode()] = &dst[i]
// 	}

// 	// Return status. WriteHeader() must be called before call to Write
// 	if r == 0 {
// 		rw.WriteHeader(http.StatusOK)
// 	} else {
// 		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
// 		return
// 	}

// 	// Return body
// 	encoder := json.NewEncoder(rw)
// 	if err = encoder.Encode(m); err != nil {
// 		log.Println(err, "in encoding result", m)
// 	} else {
// 		log.Printf("SearchBook() returns %d items\n", len(m))
// 	}
// }

// func storeTen(rw http.ResponseWriter, req *http.Request) {
// 	// Store 10 random entities
// 	r := 0
// 	c := appengine.NewContext(req)
// 	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
// 	for i := 0; i < 10; i++ {
// 		v := Book{
// 			Name:       BookName[i],
// 			Author:     BookAuthor[i],
// 			Pages:      rand.Intn(BookMaxPages),
// 			Year:       rand.Intn(time.Now().Year()),
// 			CreateTime: time.Now(),
// 		}
// 		if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, BookKind, pKey), &v); err != nil {
// 			log.Println(err)
// 			r = 1
// 			break
// 		}
// 	}

// 	// Return status. WriteHeader() must be called before call to Write
// 	if r == 0 {
// 		rw.WriteHeader(http.StatusCreated)
// 	} else {
// 		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
// 	}
// }

// func deleteBook(rw http.ResponseWriter, req *http.Request) {
// 	// Get key from URL
// 	tokens := strings.Split(req.URL.Path, "/")
// 	var keyIndexInTokens int = 0
// 	for i, v := range tokens {
// 		if v == "books" {
// 			keyIndexInTokens = i + 1
// 		}
// 	}
// 	if keyIndexInTokens >= len(tokens) {
// 		log.Println("Key is not given so that delete all books")
// 		deleteAll(rw, req)
// 		return
// 	}
// 	keyString := tokens[keyIndexInTokens]
// 	if keyString == "" {
// 		log.Println("Key is empty so that delete all books")
// 		deleteAll(rw, req)
// 	} else {
// 		deleteOneBook(rw, req, keyString)
// 	}
// }

// func deleteAll(rw http.ResponseWriter, req *http.Request) {
// 	// Delete root entity after other entities
// 	r := 0
// 	c := appengine.NewContext(req)
// 	pKey := datastore.NewKey(c, BookKind, BookRoot, 0, nil)
// 	if keys, err := datastore.NewQuery(BookKind).KeysOnly().GetAll(c, nil); err != nil {
// 		log.Println(err)
// 		r = 1
// 	} else if err := datastore.DeleteMulti(c, keys); err != nil {
// 		log.Println(err)
// 		r = 1
// 	} else if err := datastore.Delete(c, pKey); err != nil {
// 		log.Println(err)
// 		r = 1
// 	}

// 	// Return status. WriteHeader() must be called before call to Write
// 	if r == 0 {
// 		rw.WriteHeader(http.StatusOK)
// 	} else {
// 		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 	}
// }

// func deleteOneBook(rw http.ResponseWriter, req *http.Request, keyString string) {
// 	// Result
// 	r := http.StatusNoContent
// 	defer func() {
// 		// Return status. WriteHeader() must be called before call to Write
// 		if r == http.StatusNoContent {
// 			rw.WriteHeader(http.StatusNoContent)
// 		} else {
// 			http.Error(rw, http.StatusText(r), r)
// 		}
// 	}()

// 	key, err := datastore.DecodeKey(keyString)
// 	if err != nil {
// 		log.Println(err, "in decoding key string")
// 		r = http.StatusBadRequest
// 		return
// 	}

// 	// Delete the entity
// 	c := appengine.NewContext(req)
// 	if err := datastore.Delete(c, key); err != nil {
// 		log.Println(err, "in deleting entity by key")
// 		r = http.StatusNotFound
// 		return
// 	}
// 	log.Println(key, "is deleted")
// }
