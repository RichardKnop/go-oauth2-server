package gzip

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	gzipTestString              = "Foobar Wibble Content"
	gzipTestWebSocketKey        = "Test"
	gzipInvalidCompressionLevel = 11
)

func testHTTPContent(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, gzipTestString)
}

func Test_ServeHTTP_Compressed(t *testing.T) {
	gzipHandler := Gzip(DefaultCompression)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(headerAcceptEncoding, encodingGzip)

	gzipHandler.ServeHTTP(w, req, testHTTPContent)

	gr, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer gr.Close()

	body, _ := ioutil.ReadAll(gr)

	if string(body) != gzipTestString {
		t.Fail()
	}
}

func Test_ServeHTTP_NoCompression(t *testing.T) {
	gzipHandler := Gzip(NoCompression)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	gzipHandler.ServeHTTP(w, req, testHTTPContent)

	if w.Body.String() != gzipTestString {
		t.Fail()
	}
}

func Test_ServeHTTP_CompressionWithNoGzipHeader(t *testing.T) {
	gzipHandler := Gzip(DefaultCompression)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	gzipHandler.ServeHTTP(w, req, testHTTPContent)

	if w.Body.String() != gzipTestString {
		t.Fail()
	}
}

func Test_ServeHTTP_CompressionWithPrecompressedResponse(t *testing.T) {
	gzipHandler := Gzip(DefaultCompression)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAcceptEncoding, encodingGzip)

	gzipHandler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerContentEncoding, "deflate")
		testHTTPContent(w, r)
	})

	if w.Body.String() != gzipTestString {
		t.Fail()
	}
}

func Test_ServeHTTP_InvalidCompressionLevel(t *testing.T) {
	gzipHandler := Gzip(gzipInvalidCompressionLevel)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(headerAcceptEncoding, encodingGzip)

	stop := make(chan struct{})
	defer func() {
		if r := recover(); r != nil {
			close(stop)
		}
	}()

	gzipHandler.ServeHTTP(w, req, testHTTPContent)

	select {
	case <-stop:
	case <-time.After(50 * time.Millisecond):
		t.Fail()
	}

}

func Test_ServeHTTP_WebSocketConnection(t *testing.T) {
	gzipHandler := Gzip(DefaultCompression)
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(headerAcceptEncoding, encodingGzip)
	req.Header.Set(headerSecWebSocketKey, gzipTestWebSocketKey)

	gzipHandler.ServeHTTP(w, req, testHTTPContent)

	if w.Body.String() != gzipTestString {
		t.Fail()
	}
}

func Benchmark_ServeHTTP(b *testing.B) {

	b.StopTimer()
	b.ReportAllocs()

	gzipHandler := Gzip(DefaultCompression)
	req, err := http.NewRequest("GET", "http://localhost/foobar", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set(headerAcceptEncoding, encodingGzip)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		gzipHandler.ServeHTTP(w, req, testHTTPContent)
	}

}
