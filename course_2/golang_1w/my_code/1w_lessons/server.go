package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/icrowley/fake"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	client = http.Client{
		Timeout: time.Millisecond,
	}
	upgrader = websocket.Upgrader{
		HandshakeTimeout:  time.Hour,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       func (r *http.Request) bool { return true },
	}
	AvgSleep = 50
	Addr     = "localhost"
	Port     = 8080
)

type Timing struct {
	Count    int
	Duration time.Duration
}

type ctxTimings struct {
	sync.Mutex
	Data map[string]*Timing
}

type AccessLogger struct {
	File   string
	Logger *zap.SugaredLogger
}

type key int

const timingsKey key = 1

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	tpl := template.Must(template.ParseFiles("./tpl/index.html"))

	AccessLogOut := new(AccessLogger)
	AccessLogOut.Logger = InitAccessLogger()

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/", adminIndex)
	adminMux.HandleFunc("/admin/users/", adminUsers)
	adminMux.HandleFunc("/admin/orders/", adminOrders)

	adminHandler := adminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler)
	siteMux.HandleFunc("/panic/", panicPage)
	siteMux.HandleFunc("/remote/", resourseRequest)
	siteMux.HandleFunc("/posts/", loadPostsHandle)
	siteMux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	    tpl.Execute(w,nil)
	})
	siteMux.HandleFunc("/notifications", func (w http.ResponseWriter, r *http.Request) {
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		go sendNewMsgNotifications(connection)
	})

	siteHandler := timingMiddleware(siteMux)
	siteHandler = AccessLogOut.accessLogMiddleware(siteHandler)
	siteHandler = panicMiddleware(siteHandler)

	_ = http.ListenAndServe(":8080", siteHandler)
}

// --------- Middleware --------- //
// --------- ********** --------- //

func adminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("adminAuthMiddleware shooting at ", r.URL.String())
		next.ServeHTTP(w, r)
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("panicMiddleware", r.URL.String())
		defer func() {
			fmt.Println("there is a panic at ", r.URL.String())
			if err := recover(); err != nil {
				fmt.Println(err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx,
			timingsKey,
			&ctxTimings{
				Data: make(map[string]*Timing),
			})
		defer logContextTimings(ctx, r.URL.Path, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (al *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		al.Logger.Info(r.URL.Path,
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}

// --------- Middleware helpers --------- //
// --------- ********** ******* --------- //

func logContextTimings(ctx context.Context, path string, start time.Time) {

	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	totalReal := time.Since(start)
	buf := bytes.NewBufferString(path)
	var total time.Duration
	for timing, value := range timings.Data {
		total += value.Duration
		buf.WriteString(fmt.Sprintf("\n\t%s(%d): %s", timing, value.Count, value.Duration))
	}
	buf.WriteString(fmt.Sprintf("\n\ttotal: %s", totalReal))
	buf.WriteString(fmt.Sprintf("\n\ttracked: %s", total))
	buf.WriteString(fmt.Sprintf("\n\tunkn: %s", totalReal-total))

	fmt.Println(buf.String())
}

func trackContextTimings(ctx context.Context, metricName string, start time.Time) {
	timings, ok := ctx.Value(timingsKey).(*ctxTimings)
	if !ok {
		return
	}
	elapsed := time.Since(start)
	// лочимся на случай конкурентной записи в мапку
	timings.Lock()
	defer timings.Unlock()
	// если метрики ещё нет - мы её создадим, если есть - допишем в существующую
	if metric, metricExist := timings.Data[metricName]; !metricExist {
		timings.Data[metricName] = &Timing{
			Count:    1,
			Duration: elapsed,
		}
	} else {
		metric.Count++
		metric.Duration += elapsed
	}
}

// --------- Handlers --------- //
// --------- ******** --------- //

func adminIndex(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "admin Index handler")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func adminUsers(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "admin Users handler")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func adminOrders(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "admin Orders handler")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func panicPage(w http.ResponseWriter, r *http.Request) {
	panic("Aaaaaaaaaaaaaaaaaaaaaa!")
}

func resourseRequest(w http.ResponseWriter, r *http.Request) {
	err := getRemoteResource()
	if err != nil {
		switch err := errors.Cause(err).(type) {
		case *url.Error:
			fmt.Printf("resourse %s error: %+v\n", err.URL, err.Err)
			http.Error(w, "remote resource error", 500)
		default:
			fmt.Printf("%+v\n", err)
			http.Error(w, "parsing error", 500)
		}
		return
	}

	_, err = w.Write([]byte("it's going OK"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// --------- Functions --------- //
// --------- ********* --------- //

func getRemoteResource() error {
	url := "http://127.0.0.1:9999/pages?id=123"
	_, err := client.Get(url)
	if err != nil {
		return errors.Wrap(err, "resource error")
	}
	return nil
}

func loadPostsHandle(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	emulateWork(ctx, "checkCache")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	emulateWork(ctx, "loadPosts")
	time.Sleep(10 * time.Millisecond)
	emulateWork(ctx, "loadSidebar")
	emulateWork(ctx, "loadComments")

	fmt.Fprintln(w, "Request done")
}

func emulateWork(ctx context.Context, workName string) {
	defer trackContextTimings(ctx, workName, time.Now())

	rnd := time.Duration(rand.Intn(AvgSleep))
	time.Sleep(time.Millisecond * rnd)
}

func InitAccessLogger() *zap.SugaredLogger {
	writerSyncer := getAccessLogWriter()
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger.Sugar()
}

func getAccessLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create("./test.log")
	return zapcore.AddSync(file)
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   fake.EmailAddress(),
		"name":    fake.FirstName() + " " + fake.LastName(),
		"subject": fake.Product() + " " + fake.Model(),
	})
	return data
}

func sendNewMsgNotifications(client *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage) // get writer for next message
		if err != nil {
			ticker.Stop()
			break
		}

		msg := newMessage() // get the message
		w.Write(msg) // write message to writer
		w.Close() // close writer

		<-ticker.C
	}
}
