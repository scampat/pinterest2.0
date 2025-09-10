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

// Modelo simplificado de Pin (lo puedes expandir después)
type Pin struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	BoardID string `json:"board_id"`
	Title   string `json:"title"`
	URL     string `json:"url"` // URL de la imagen
}

// Gateway permite que users hable con pins
type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetPinsByUser obtiene todos los pins de un usuario desde el servicio pins
func (g *Gateway) GetPinsByUser(ctx context.Context, userID string) ([]Pin, error) {
	// Pedimos al registry las direcciones de pins
	addrs, err := g.registry.ServiceAddress(ctx, "pins")
	if err != nil {
		return nil, err
	}

	// Seleccionamos una al azar (load balancing simple)
	url := "http://" + addrs[rand.Intn(len(addrs))] + "/pins?user_id=" + userID
	log.Printf("Calling pins service: GET %s", url)

	// Construimos request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// Ejecutamos request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Manejo de errores HTTP
	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp.Status)
	}

	// Decodificamos respuesta JSON
	var pins []Pin
	if err := json.NewDecoder(resp.Body).Decode(&pins); err != nil {
		return nil, err
	}

	return pins, nil
}
