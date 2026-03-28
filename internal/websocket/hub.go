package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients by company_id and branch_id
	clients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data"`
	CompanyID uuid.UUID   `json:"company_id,omitempty"`
	BranchID  uuid.UUID   `json:"branch_id,omitempty"`
}

var instance *Hub
var once sync.Once

// GetHub returns the singleton Hub instance
func GetHub() *Hub {
	once.Do(func() {
		instance = &Hub{
			broadcast:  make(chan *Message, 256),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			clients:    make(map[string]map[*Client]bool),
		}
		go instance.run()
	})
	return instance
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			key := client.GetKey()
			if _, ok := h.clients[key]; !ok {
				h.clients[key] = make(map[*Client]bool)
			}
			h.clients[key][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (Total clients for key: %d)", key, len(h.clients[key]))

		case client := <-h.unregister:
			h.mu.Lock()
			key := client.GetKey()
			if clients, ok := h.clients[key]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, key)
					}
					log.Printf("Client unregistered: %s (Remaining clients for key: %d)", key, len(clients))
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			key := getKey(message.CompanyID, message.BranchID)
			clients := h.clients[key]
			h.mu.RUnlock()

			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			for client := range clients {
				select {
				case client.send <- messageBytes:
				default:
					h.mu.Lock()
					close(client.send)
					delete(h.clients[key], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

// BroadcastOrderUpdate broadcasts order update to all clients in the same company/branch
func (h *Hub) BroadcastOrderUpdate(action string, data interface{}, companyID, branchID uuid.UUID) {
	message := &Message{
		Type:      "order",
		Action:    action,
		Data:      data,
		CompanyID: companyID,
		BranchID:  branchID,
	}
	h.broadcast <- message
}

func getKey(companyID, branchID uuid.UUID) string {
	return companyID.String() + ":" + branchID.String()
}
