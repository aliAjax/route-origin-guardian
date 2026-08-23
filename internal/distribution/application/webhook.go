package application

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"net/http"
	"time"
)

type WebhookSender struct {
	Client      *http.Client
	URL, Secret string
	Attempts    int
}

func (w WebhookSender) Send(ctx context.Context, e domain.Event) error {
	if w.Client == nil {
		w.Client = &http.Client{Timeout: 5 * time.Second}
	}
	b, x := json.Marshal(e)
	if x != nil {
		return x
	}
	m := hmac.New(sha256.New, []byte(w.Secret))
	_, _ = m.Write(b)
	req, x := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytes.NewReader(b))
	if x != nil {
		return x
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Event-ID", fmt.Sprint(e.ID))
	req.Header.Set("X-Signature", hex.EncodeToString(m.Sum(nil)))
	resp, x := w.Client.Do(req)
	if x != nil {
		return fmt.Errorf("webhook: %w", x)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook status %s", resp.Status)
	}
	return nil
}
