package handler

//
//import (
//	"github.com/bbt-t/shortenerURL/internal/controller/rest"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/bbt-t/shortenerURL/internal/adapter/storage"
//	"github.com/bbt-t/shortenerURL/internal/config"
//)
//
//var cfg = &configs.ServerCfg{
//	ServerAddress: "127.0.0.1:5001",
//	BaseURL:       "http://127.0.0.1:5001",
//}
//
//type Fields struct {
//	urls map[string]string
//}
//
//func (f Fields) TestHandler_takeAndSendURL(t *testing.T) {
//	type want struct {
//		code int
//		url  string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "test POST #1",
//			want: want{
//				code: 201,
//				url:  "https://ya.ru",
//			},
//		},
//	}
//	for _, tt := range tests {
//		for t.Run(tt.name, func(t *testing.T) {}) {
//			request := httptest.NewRequest(http.MethodPost, "/", nil)
//
//			w := httptest.NewRecorder()
//			db := storage.NewMapDBPlug()
//			th := rest.NewHttpServer(db, *cfg)
//
//			appH := http.HandlerFunc(th.composeNewShortURL)
//			appH.ServeHTTP(w, request)
//			res := w.Result()
//			if res.StatusCode != tt.want.code {
//				t.Errorf("Exepted status code %d, got %d", tt.want.code, w.Code)
//			}
//			defer res.Body.Close()
//			resBody, err := io.ReadAll(res.Body)
//			if err != nil {
//				t.Fatal()
//			}
//			shortURL := string(resBody)[len("http://localhost:8080/"):]
//			f.urls[tt.want.url] = shortURL
//		}
//	}
//}
//
//func (f Fields) TestHandler_redirectToOriginalURL(t *testing.T) {
//	type want struct {
//		code     int
//		location string
//	}
//	tests := []struct {
//		name string
//		want want
//	}{
//		{
//			name: "test GET #1",
//			want: want{
//				code:     307,
//				location: "https://ya.ru",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			request := httptest.NewRequest(http.MethodGet, "/"+tt.want.location, nil)
//			w := httptest.NewRecorder()
//			db := storage.NewMapDBPlug()
//			th := NewHandlerServer(db, *cfg)
//			appH := http.HandlerFunc(th.recoverOriginalURL)
//			appH.ServeHTTP(w, request)
//			res := w.Result()
//			defer res.Body.Close()
//			if res.StatusCode != tt.want.code {
//				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
//			}
//			resLocation := res.Header.Get("Location")
//			if resLocation != tt.want.location {
//				t.Errorf("Expected location %s, got %s", tt.want.location, resLocation)
//			}
//		})
//	}
//}
