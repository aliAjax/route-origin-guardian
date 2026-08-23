package httpx

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/platform/config"
	"github.com/routeorigin/route-origin-guardian/internal/rib/application"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	rpkidomain "github.com/routeorigin/route-origin-guardian/internal/rpki/domain"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cfg  config.Config
	svc  *application.Service
	roas *rpkidomain.Index
	http *http.Server
}

func NewServer(cfg config.Config, svc *application.Service, roas *rpkidomain.Index) *Server {
	s := &Server{cfg: cfg, svc: svc, roas: roas}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.metrics)
	mux.HandleFunc("/v1/routes", s.routes)
	mux.HandleFunc("/v1/events", s.events)
	mux.HandleFunc("/v1/rpki/validate", s.validate)
	mux.HandleFunc("/v1/rpki/snapshots", s.snapshot)
	s.http = &http.Server{Addr: cfg.HTTPListen, Handler: s.middleware(mux), ReadTimeout: cfg.ReadTimeout, WriteTimeout: cfg.WriteTimeout}
	return s
}
func (s *Server) ListenAndServe() error {
	slog.Info("http listening", "addr", s.cfg.HTTPListen)
	return s.http.ListenAndServe()
}
func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			if x := recover(); x != nil {
				writeErr(w, http.StatusInternalServerError, "internal error")
			}
			slog.Info("http request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start).String())
		}()
		if r.URL.Path != "/healthz" && r.URL.Path != "/readyz" && r.URL.Path != "/metrics" {
			key := r.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(key), []byte(s.cfg.APIKey)) != 1 {
				writeErr(w, http.StatusUnauthorized, "missing or invalid API key")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if e := s.svc.Ready(r.Context()); e != nil {
		writeErr(w, 503, e.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (s *Server) metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# HELP route_origin_guardian_info Service information.")
	fmt.Fprintln(w, "route_origin_guardian_info 1")
}
func (s *Server) routes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := application.Query{TenantID: r.URL.Query().Get("tenant_id"), Prefix: r.URL.Query().Get("prefix"), Peer: r.URL.Query().Get("peer"), Family: r.URL.Query().Get("family")}
		if n, _ := strconv.Atoi(r.URL.Query().Get("asn")); n > 0 {
			q.ASN = uint32(n)
		}
		rows, e := s.svc.List(r.Context(), q)
		if e != nil {
			writeErr(w, 500, e.Error())
			return
		}
		writeJSON(w, 200, rows)
	case http.MethodPost:
		var rawBody json.RawMessage
		if e := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&rawBody); e != nil {
			writeErr(w, 400, e.Error())
			return
		}
		var in domain.Route
		if e := json.Unmarshal(rawBody, &in); e != nil {
			writeErr(w, 400, e.Error())
			return
		}
		if in.TenantID == "" {
			var envelope struct {
				Operation string       `json:"operation"`
				Route     domain.Route `json:"route"`
			}
			if e := json.Unmarshal(rawBody, &envelope); e == nil && envelope.Route.TenantID != "" {
				in = envelope.Route
				in.Operation = domain.Operation(envelope.Operation)
			}
		}
		if in.Operation == domain.Withdraw {
			e := func() error { _, x := s.svc.Withdraw(r.Context(), in); return x }()
			if e != nil {
				writeErr(w, 422, e.Error())
				return
			}
			writeJSON(w, 202, map[string]string{"status": "withdrawn"})
			return
		}
		ev, e := s.svc.Announce(r.Context(), in)
		if e != nil {
			writeErr(w, 422, e.Error())
			return
		}
		writeJSON(w, 202, ev)
	default:
		w.WriteHeader(405)
	}
}
func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	after, _ := strconv.ParseUint(r.URL.Query().Get("after"), 10, 64)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	rows, e := s.svc.Events(r.Context(), after, limit)
	if e != nil {
		writeErr(w, 500, e.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"items": rows, "next": func() uint64 {
		if len(rows) == 0 {
			return after
		}
		return rows[len(rows)-1].ID
	}()})
}
func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("prefix")
	asn, _ := strconv.ParseUint(r.URL.Query().Get("asn"), 10, 32)
	if p == "" || asn == 0 {
		writeErr(w, 400, "prefix and asn required")
		return
	}
	writeJSON(w, 200, s.roas.Validate(p, uint32(asn)))
}
func (s *Server) snapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeErr(w, 405, "method not allowed")
		return
	}
	var body struct {
		ROAs []rpkidomain.ROA `json:"roas"`
	}
	if e := json.NewDecoder(r.Body).Decode(&body); e != nil {
		writeErr(w, 400, e.Error())
		return
	}
	if e := s.roas.Replace(body.ROAs); e != nil {
		writeErr(w, 422, e.Error())
		return
	}
	writeJSON(w, 202, map[string]any{"version": s.roas.Version(), "count": len(body.ROAs)})
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

var _ = strings.TrimSpace
