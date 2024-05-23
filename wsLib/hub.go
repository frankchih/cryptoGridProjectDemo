package wsLib

type BroadcastChannel struct {
	Channel string
	Message []byte
}

type Hub struct {
	// 存儲所有的 channel 和 client
	Channels map[string]map[*Client]bool

	// 用來新增或刪除 channel 和 client
	Register   chan *Client
	Unregister chan *Client

	// 用來傳遞訊息
	Broadcast        chan []byte
	BroadcastChannel chan BroadcastChannel
}

func NewHub() *Hub {
	return &Hub{
		Channels:         make(map[string]map[*Client]bool),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Broadcast:        make(chan []byte),
		BroadcastChannel: make(chan BroadcastChannel),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			clients, ok := h.Channels[client.Channel]
			if !ok {
				clients = make(map[*Client]bool)
				h.Channels[client.Channel] = clients
			}
			clients[client] = true
		case client := <-h.Unregister:
			clients, ok := h.Channels[client.Channel]
			if ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Channels, client.Channel)
					}
				}
			}
		case broadcastChannel := <-h.BroadcastChannel:
			// 將訊息廣播給指定的 channel 的所有 client
			//fmt.Println(h.channels)
			channel := broadcastChannel.Channel
			message := broadcastChannel.Message
			clients, ok := h.Channels[channel]
			//fmt.Println("broadcastChannel", len(message), len(h.BroadcastChannel))
			if ok {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(h.Channels, client.Channel)
						}
					}
				}
			}
			//case message := <-h.Broadcast:
			//	// 將訊息廣播給指定的 channel 的所有 client
			//	clients, ok := h.Channels["ws/quote/"]
			//	fmt.Println(clients, ok)
			//	if ok {
			//		for client := range clients {
			//			select {
			//			case client.Send <- message:
			//			default:
			//				close(client.Send)
			//				delete(clients, client)
			//				if len(clients) == 0 {
			//					delete(h.Channels, client.Channel)
			//				}
			//			}
			//		}
			//	}

		}

	}
}
