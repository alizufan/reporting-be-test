package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"reporting/db"
	"reporting/libs/logger"
	"reporting/libs/util"

	AuthHandler "reporting/handler/auth"
	TrxHandler "reporting/handler/transaction"
	MerchantRepo "reporting/repository/merchant"
	TrxRepo "reporting/repository/transaction"
	UserRepo "reporting/repository/user"
	AuthSrv "reporting/service/auth"
	TrxSrv "reporting/service/transaction"

	AppMiddleware "reporting/server/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-rel/rel"
)

func NewHTTPServer() *HTTPServer {
	d := db.Init()
	dbRepo := rel.New(d)

	logger.New()
	util.NewValidator()

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		ExposedHeaders:     []string{"*"},
		AllowCredentials:   true,
		MaxAge:             60,
		OptionsPassthrough: false,
		Debug:              false,
	}))
	r.Use(AppMiddleware.Tracker)

	UserR := UserRepo.User{
		DB: dbRepo,
	}
	MerchantR := MerchantRepo.Merchant{
		DB: dbRepo,
	}
	AuthS := AuthSrv.Auth{
		UserRepo:     &UserR,
		MerchantRepo: &MerchantR,
	}
	AuthH := AuthHandler.AuthHandler{
		AuthService: &AuthS,
	}

	TrxR := TrxRepo.Transaction{
		DB: dbRepo,
	}
	TrxS := TrxSrv.Transaction{
		TransactionRepo: &TrxR,
	}
	TrxH := TrxHandler.TransactionHandler{
		TransactionSrv: &TrxS,
	}

	server := &HTTPServer{
		Router:            r,
		DB:                d,
		TransactionRouter: &TrxH,
		AuthRouter:        &AuthH,
	}

	server.routes()

	return server
}

type HTTPServer struct {
	Router *chi.Mux
	DB     rel.Adapter

	AuthRouter        AuthRouter
	TransactionRouter TransactionRouter
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	port, ok := os.LookupEnv("API_PORT")
	if !ok {
		port = "3000"
	}

	server := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           hs.Router,
		IdleTimeout:       0,
		WriteTimeout:      5 * time.Second,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		log.Printf("start reporting api")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start / shutdown reporting api, err : \n%+v\n", err)
		}
	}()

	<-ctx.Done()

	shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("server shutdown")

	if err := server.Shutdown(shutdown); err != nil {
		log.Fatalf("shutdown reporting api, err : \n%+v\n", err)
	}

	log.Printf("server shutdown properly")

	if err := hs.DB.Close(); err != nil {
		log.Fatal("unable close db connection")
	}

	return nil
}

type AuthRouter interface {
	Login(rw http.ResponseWriter, r *http.Request)
}

type TransactionRouter interface {
	Report(rw http.ResponseWriter, r *http.Request)
}

func (hs *HTTPServer) routes() {
	hs.Router.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("reporting api 🔥"))
	})

	hs.Router.Post("/login", hs.AuthRouter.Login)
	hs.Router.With(AppMiddleware.JWTValidation).Get("/report", hs.TransactionRouter.Report)
}
