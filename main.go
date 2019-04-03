/*
 * hub-kubernetes-agent
 *
 * an agent used to provision and configure Kubernetes resources
 *
 * API version: v1beta
 * Contact: support@appvia.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/urfave/cli"

	sw "github.com/appvia/hub-kubernetes-agent/go"
	muxlogrus "github.com/pytimer/mux-logrus"
	logrus "github.com/sirupsen/logrus"
)

var (
	release = "v0.0.1"
)

func invokeServerAction(ctx *cli.Context) error {
	router := sw.NewRouter()
	router.Use(Middleware)

	var logoptions muxlogrus.LogOptions
	logoptions = muxlogrus.LogOptions{Formatter: new(logrus.JSONFormatter), EnableStarting: true}
	router.Use(muxlogrus.NewLogger(logoptions).Middleware)

	srv := &http.Server{
		Addr:         ctx.String("listen") + ":" + ctx.String("http-port"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("failed to start the api service")
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChannel

	return nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "healthz") == true {
			next.ServeHTTP(w, r)
			return
		}
		tokenHeader := r.Header.Get("Authorization")
		if len(tokenHeader) == 0 {
			logrus.Infof("Missing Authorization header")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		token := strings.Replace(tokenHeader, "Bearer ", "", 1)
		if token != os.Getenv("AUTH_TOKEN") {
			logrus.Infof("Incorrect Authorization header")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if r.Header.Get("X-Kube-API-URL") == "" || r.Header.Get("X-Kube-Token") == "" || r.Header.Get("X-Kube-CA") == "" {
			logrus.Infof("Missing Kube header")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	app := &cli.App{
		Name:    "hub-kubernetes-agent",
		Author:  "Daniel Whatmuff",
		Email:   "daniel.whatmuff@appvia.io",
		Usage:   "A backend agent used to provision resources within Kubernetes clusters",
		Version: release,

		OnUsageError: func(context *cli.Context, err error, _ bool) error {
			fmt.Fprintf(os.Stderr, "[error] invalid options %s\n", err)
			return err
		},

		Action: func(ctx *cli.Context) error {
			if ctx.String("auth-token") == "" {
				return cli.NewExitError("Missing AUTH_TOKEN", 1)
			}
			os.Setenv("AUTH_TOKEN", ctx.String("auth-token"))
			logrus.Info("Starting server...")
			return invokeServerAction(ctx)
		},

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "listen",
				Usage:  "the interface to bind the service to `INTERFACE`",
				Value:  "127.0.0.1",
				EnvVar: "LISTEN",
			},
			cli.IntFlag{
				Name:   "http-port",
				Usage:  "network interface the service should listen on `PORT`",
				Value:  10080,
				EnvVar: "HTTP_PORT",
			},
			cli.StringFlag{
				Name:   "auth-token",
				Usage:  "authentication token used to verifier the caller `TOKEN`",
				EnvVar: "AUTH_TOKEN",
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}
