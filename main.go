package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	// "github.com/bradleyfalzon/ghinstallation"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lmittmann/tint"
)

const defaultPort = "8080"

var (
	githubAppId   int64
	webhookSecret string
	rsaPrivateKey *rsa.PrivateKey
)

var roundTripper = http.DefaultTransport

func main() {
	slog.SetDefault(setUpLogging())

	var err error
	githubAppId, err = strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if err != nil {
		slog.Error("APP_ID env var not set or not a valid int64", slog.Any("error", err))
		os.Exit(1)
	}
	webhookSecret = os.Getenv("GITHUB_APP_WEBHOOK_SECRET")
	if webhookSecret == "" {
		slog.Error("WEBHOOK_SECRET env var not set")
		os.Exit(1)
	}
	privateKeyBase64 := os.Getenv("GITHUB_APP_PRIVATE_KEY_BASE64")
	if privateKeyBase64 == "" {
		slog.Error("GITHUB_APP_PRIVATE_KEY_BASE64 env var not set")
		os.Exit(1)
	}
	privateKey, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		slog.Error("error decoding GitHub App private key from base64", slog.Any("error", err))
		os.Exit(1)
	}

	rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		slog.Error("error parsing GitHub App RSA private key from PEM", slog.Any("error", err))
	}

	slog.Info("server is starting")
	port := os.Getenv("PORT")
	if port == "" {
		slog.Info(fmt.Sprintf("PORT env var not set, using default port %s", defaultPort))
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", index)
	mux.Handle("POST /webhook", WithWebhookSecret(
		WithAuthenticatedApp( // provides gh_app_client
			WithAuthenticatedAppInstallation( // provides gh_installation_client
				http.HandlerFunc(handleWebhook),
			),
		),
	),
	)

	err = http.ListenAndServe(fmt.Sprint("0.0.0.0:", port), mux)
	if err != nil {
		slog.Error("failed to start listening", slog.Any("error", err))
		os.Exit(1)
	}
}

func WithWebhookSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtain the signature from the request
		theirSignature := r.Header.Get("X-Hub-Signature-256")
		parts := strings.Split(theirSignature, "=")
		if len(parts) != 2 {
			http.Error(w, "invalid webhook signature", http.StatusForbidden)
			return
		}
		theirHexMac := parts[1]
		theirMac, err := hex.DecodeString(theirHexMac)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding webhook signature: %v", err), http.StatusBadRequest)
			return
		}

		// Calculate our own signature
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(payload)) // make body available for reading again

		hash := hmac.New(sha256.New, []byte(webhookSecret))
		hash.Write(payload)
		ourMac := hash.Sum(nil)

		// Compare signatures
		if !hmac.Equal(theirMac, ourMac) {
			http.Error(w, "webhook signature is invalid", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func WithAuthenticatedApp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := jwt.MapClaims{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(10 * time.Minute).Unix(),
			"iss": githubAppId,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenStr, err := token.SignedString(rsaPrivateKey)
		if err != nil {
			slog.Error("error signing JWT", slog.Any("error", err))
			http.Error(w, "error signing JWT", http.StatusInternalServerError)
			return
		}

		appClient := http.Client{Transport: &BearerTransport{Token: tokenStr}}

		ctx := context.WithValue(r.Context(), "gh_app_client", appClient)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}

func WithAuthenticatedAppInstallation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read request body
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		body := bytes.NewBuffer(b)
		r.Body = io.NopCloser(bytes.NewBuffer(bytes.Clone(b))) // make body available for reading again

		// decode body from JSON into a map
		decoder := json.NewDecoder(body)
		decoder.UseNumber()

		var payload map[string]interface{}
		err = decoder.Decode(&payload)
		if err != nil {
			slog.Error("error reading body", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error reading body: %v", err)
			return
		}

		// extract installation ID from the request body

		installationIdStr := payload["installation"].(map[string]interface{})["id"].(json.Number).String()
		installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
		if err != nil {
			slog.Error("error parsing installation id", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Error parsing installation id: %v", err)
			return
		}

		// get app installation access token

		url := fmt.Sprintf("https://api.github.com/app/installations/%d/access_tokens", installationId)
		appClient := r.Context().Value("gh_app_client").(http.Client)
		res, err := appClient.Post(url, "application/json", nil)
		if err != nil {
			msg := "error calling endpoint " + url
			slog.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// read response body and extract the access token

		// read body
		b, err = io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
			return
		}
		body = bytes.NewBuffer(b)

		// decode body from JSON into a map
		decoder = json.NewDecoder(body)
		decoder.UseNumber()

		clear(payload)
		err = decoder.Decode(&payload)
		if err != nil {
			msg := fmt.Sprint("error reading response body from " + res.Request.URL.String())
			slog.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		appInstallationToken := payload["token"].(string)
		appInstallationClient := http.Client{
			Transport: &BearerTransport{Token: appInstallationToken},
		}
		slog.Info("installation access token obtained", slog.Any("token", appInstallationToken))

		ctx := context.WithValue(r.Context(), "gh_installation_client", appInstallationClient)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}

func setUpLogging() *slog.Logger {
	// Configure logging
	logLevel := slog.LevelDebug
	prod := os.Getenv("K_SERVICE") != "" // https://cloud.google.com/run/docs/container-contract#services-env-vars
	if prod {
		// Based on https://github.com/remko/cloudrun-slog
		const LevelCritical = slog.Level(12)
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     logLevel,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.MessageKey:
					a.Key = "message"
				case slog.SourceKey:
					a.Key = "logging.googleapis.com/sourceLocation"
				case slog.LevelKey:
					a.Key = "severity"
					level := a.Value.Any().(slog.Level)
					if level == LevelCritical {
						a.Value = slog.StringValue("CRITICAL")
					}
				}
				return a
			},
		}

		gcpHandler := slog.NewJSONHandler(os.Stderr, opts)
		return slog.New(gcpHandler)
	} else {
		opts := tint.Options{Level: logLevel, TimeFormat: time.TimeOnly, AddSource: true}
		handler := tint.NewHandler(os.Stdout, &opts)
		return slog.New(handler)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	slog.Info("request received", slog.String("path", r.URL.Path))
	fmt.Fprintln(w, "hello world")
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	l := slog.With(slog.String("path", r.URL.Path))

	eventType := r.Header.Get("X-GitHub-Event")

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	var payload map[string]interface{}
	err := decoder.Decode(&payload)
	if err != nil {
		msg := "error reading request body"
		l.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	installationIdStr := payload["installation"].(map[string]interface{})["id"].(json.Number).String()
	installationId, err := strconv.ParseInt(installationIdStr, 10, 64)
	if err != nil {
		l.Error("error parsing installation id", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing installation id: %v", err)
		return
	}

	action, _ := payload["action"].(string)
	l.Info("new request",
		slog.String("event", eventType),
		slog.Any("action", action),
		slog.Int64("id", installationId),
	)

	// transport, err := ghinstallation.New(roundTripper, githubAppId, installationId, privateKey)
	// if err != nil {
	// 	l.Error("error creating transport", slog.Any("error", err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprintf(w, "Error creating transport: %v", err)
	// 	return
	// }
	// _ = transport

	// Check type of webhook event
	switch eventType {
	case "installation":
		installation := payload["installation"].(map[string]interface{})
		login := installation["account"].(map[string]interface{})["login"].(string)

		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#installation
		if payload["action"] == "created" {
			repositories := payload["repositories"].([]interface{})
			l.Info("app installation created", slog.Any("id", installation["id"]), slog.String("login", login), slog.Int("repositories", len(repositories)))
		}
		if payload["action"] == "deleted" {
			l.Info("app installation deleted", slog.Any("id", installation["id"]), slog.String("login", login))
		}
	case "check_suite":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=requested#check_suite
		if payload["action"] == "requested" || payload["action"] == "rerequested" {
			repository := payload["repository"].(map[string]interface{})
			repoName := repository["name"].(string)
			repoOwner := repository["owner"].(map[string]interface{})["login"].(string)

			checkSuite := payload["check_suite"].(map[string]interface{})
			headSHA := checkSuite["head_sha"].(string)
			l.Info("check suite requested", slog.String("owner", repoOwner), slog.String("repo", repoName), slog.String("head_sha", headSHA))

			err := createCheckRun(r.Context(), repoOwner, repoName, headSHA)
			if err != nil {
				l.Error("error creating check run", slog.Any("error", err))
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Error creating check run: %v", err)
			}
		}
	case "check_run":
		// https://docs.github.com/en/webhooks/webhook-events-and-payloads?actionType=created#check_run
		repository := payload["repository"].(map[string]interface{})
		repoName := repository["name"].(string)
		repoOwner := repository["owner"].(map[string]interface{})["login"].(string)

		checkRun := payload["check_run"].(map[string]interface{})
		headSHA := checkRun["head_sha"].(string)
		err := createCheckRun(r.Context(), repoOwner, repoName, headSHA)
		if err != nil {
			msg := "error creating check run"
			l.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	default:
		l.Error("unknown event type", slog.String("type", eventType))
	}
}

// TODO: accept context, and access logger and authenticated HTTP client from there?

func createCheckRun(ctx context.Context, owner, repo, sha string) error {
	body := map[string]interface{}{
		"name":        "hello from bartek",
		"head_sha":    sha,
		"details_url": "https://garden.pacia.com",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling body: %w", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/check-runs", owner, repo)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	slog.Info("request created", slog.String("url", url), slog.String("body", string(bodyBytes)))

	req.Header.Set("Accept", "application/vnd.github+json")

	githubInstallationClient := ctx.Value("gh_installation_client").(http.Client)

	res, err := githubInstallationClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	respBody := make([]byte, 0)
	_, err = res.Body.Read(respBody)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	slog.Info("request made", slog.Int("status", res.StatusCode), slog.String("body", string(respBody)))

	return nil
}
