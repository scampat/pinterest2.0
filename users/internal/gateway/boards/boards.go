package boards

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	discovery "pinterest2.0/pkg/registry"
	"pinterest2.0/users/internal/gateway"
)

type Board struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Title  string `json:"title"`
}

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Consultar boards por usuario
func (g *Gateway) GetBoardsByUser(ctx context.Context, userID string) ([]Board, error) {
	addrs, err := g.registry.ServiceAddress(ctx, "boards")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/boards?user_id=" + userID
	log.Printf("Calling boards service: GET %s", url)

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

	var boards []Board
	if err := json.NewDecoder(resp.Body).Decode(&boards); err != nil {
		return nil, err
	}

	return boards, nil
}

// Crear un board desde users
func (g *Gateway) CreateBoard(ctx context.Context, userID, title string) (*Board, error) {
	addrs, err := g.registry.ServiceAddress(ctx, "boards")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/boards"
	log.Printf("Calling boards service: POST %s", url)

	payload := Board{
		UserID: userID,
		Title:  title,
	}

	data, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp.Status)
	}

	var b Board
	if err := json.NewDecoder(resp.Body).Decode(&b); err != nil {
		return nil, err
	}

	return &b, nil
}
