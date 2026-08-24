package application

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	deliverydomain "github.com/routeorigin/route-origin-guardian/internal/distribution/domain"
	"github.com/routeorigin/route-origin-guardian/internal/rib/domain"
	"io"
	"net/http"
	"time"
)

type WebhookSender struct {
	Client      *http.Client
	URL, Secret string
	Attempts    int
}

type DeadLetterStore interface {
	Put(context.Context, domain.Event) error
}

type DeliveryRunner struct {
	Sender      WebhookSender
	DeadLetter  DeadLetterStore
	MaxAttempts int
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
	_, drainErr := io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 300 {
		statusErr := fmt.Errorf("webhook status %s", resp.Status)
		return errors.Join(statusErr, drainErr)
	}
	return drainErr
}

func (r DeliveryRunner) Run(ctx context.Context, delivery *deliverydomain.Delivery) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if delivery == nil || !delivery.Begin() {
		return errors.New("delivery is not pending")
	}
	err := r.Sender.Send(ctx, delivery.Event)
	if err == nil {
		if !delivery.Complete() {
			return errors.New("delivery completion rejected")
		}
		return nil
	}
	if !delivery.Fail(r.MaxAttempts) {
		return errors.Join(err, errors.New("delivery failure transition rejected"))
	}
	if delivery.Status == deliverydomain.DeadLetter {
		if r.DeadLetter == nil {
			return errors.Join(err, errors.New("dead-letter store is nil"))
		}
		if storeErr := r.DeadLetter.Put(ctx, delivery.Event); storeErr != nil {
			return errors.Join(err, storeErr)
		}
	}
	return err
}
