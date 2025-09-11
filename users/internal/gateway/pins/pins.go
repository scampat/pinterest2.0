package pins

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	discovery "pinterest2.0/pkg/registry"
	"pinterest2.0/users/internal/gateway"
)

type Pin struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetPinsByUser(ctx context.Context, userID string) ([]Pin, error) {
	// Pedimos al registry las direcciones de pins
	addrs, err := g.registry.ServiceAddress(ctx, "pins")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/pins?user_id=" + userID
	log.Printf("Calling pins service: GET %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp.Status)
	}

	var pins []Pin
	if err := json.NewDecoder(resp.Body).Decode(&pins); err != nil {
		return nil, err
	}

	return pins, nil
}
