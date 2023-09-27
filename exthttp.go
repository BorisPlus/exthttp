package exthttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
)

type Logger interface {
	Info(msg string, a ...any)
	Error(msg string, a ...any)
}

var (
	log  *Log
	once sync.Once
)

type Log struct {
	logger Logger
}

func (l Log) Info(msg string, a ...any) {
	l.logger.Info(fmt.Sprintf("%s %s", "ExtHTTP", msg), a...)
}

func (l Log) Error(msg string, a ...any) {
	l.logger.Error(fmt.Sprintf("%s %s", "ExtHTTP", msg), a...)
}

func logInstance() *Log {
	if log == nil {
		panic("Log was not init by `initialize(logger interfaces.Logger)`.")
	}
	return log
}

func initialize(logger Logger) *Log {
	once.Do(func() {
		log = &Log{}
		log.logger = logger
	})
	return log
}

type InternalTestHTTPServer struct {
	server *http.Server
}

func NewInternalTestHTTPServer(
	host string,
	port uint16,
	logger Logger,
	headersLogsDir string,
) *InternalTestHTTPServer {
	initialize(logger)
	mux := http.NewServeMux()
	mux.Handle("/image.jpg", ImageResponser{headersLogsDir: headersLogsDir})
	mux.Handle("/", http.HandlerFunc(handleText))
	server := http.Server{
		Addr:    net.JoinHostPort(host, fmt.Sprint(port)),
		Handler: mux,
	}
	this := &InternalTestHTTPServer{}
	this.server = &server
	return this
}

func (s *InternalTestHTTPServer) Start() error {
	logInstance().Info("InternalTestHTTPServer.Start()")
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logInstance().Error(err.Error())
		return err
	}
	return nil
}

func (s *InternalTestHTTPServer) Stop(ctx context.Context) error {
	logInstance().Info("InternalTestHTTPServer.Stop()")
	if err := s.server.Shutdown(ctx); err != nil {
		logInstance().Error(err.Error())
		return err
	}
	return nil
}

type ImageResponser struct {
	headersLogsDir string
}

func (h ImageResponser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	green := image.NewRGBA(image.Rect(0, 0, 100, 100))
	draw.Draw(green, green.Bounds(), &image.Uniform{color.RGBA{0, 255, 0, 255}}, image.Point{}, draw.Src)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, green, nil); err != nil {
		text := []byte("JPEG encode error")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Content-Type", "text/plain")
		rw.Header().Set("Content-Length", strconv.Itoa(len(text)))
		_, err := rw.Write(text)
		if err != nil {
			logInstance().Error(err.Error())
		}
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "image/jpeg")
	rw.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	_, err := rw.Write(buffer.Bytes())
	if err != nil {
		logInstance().Error(err.Error())
	}
	if h.headersLogsDir != "" {
		path := path.Join(h.headersLogsDir, "headers.json")
		_ = os.Remove(path)
		file, err := os.Create(path)
		if err != nil {
			logInstance().Error(err.Error())
		}
		defer file.Close()
		jsonData := []byte{}
		err = json.Unmarshal(jsonData, &r.Header)
		if err != nil {
			logInstance().Error(err.Error())
		}
		_, err = file.Write(jsonData)
		if err != nil {
			logInstance().Error(err.Error())
		}
	}
}

func handleText(rw http.ResponseWriter, _ *http.Request) {
	text := []byte("I receive teapot-status code!")
	rw.WriteHeader(http.StatusTeapot)
	rw.Header().Set("Content-Type", "text/plain")
	rw.Header().Set("Content-Length", strconv.Itoa(len(text)))
	_, err := rw.Write(text)
	if err != nil {
		logInstance().Error(err.Error())
	}
}
