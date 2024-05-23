import { activityLogAction } from "../slices/activityLogSlice"
import { quoteActions } from "../slices/quoteSlice"
import { wsActions } from "../slices/wsSlice"
import { WebSocketObj } from "./webSocketObj"

const wsUrl = process.env.REACT_APP_WS_URL

export const webSocketMiddleware = (store) => {
    
    let webSocketObj
    return (next) => (action) => {
        if (wsActions.startConnecting.match(action)) {
            console.log("ws middleware in startConnecting")
            webSocketObj = new WebSocketObj()
            webSocketObj.connect(`${wsUrl}/ws/quote/`)
            // console.log("websocket", webSocketObj)
            webSocketObj.on("open", () => {
                // console.log("websocket obj open")
                store.dispatch(wsActions.connectionEstablished())
            })
            // setTimeout(() => {
            //     webSocketObj.send("web js")
            // }, 1000)
            
            webSocketObj.on("message", (event) => {
                // console.log(event.data.length, event.data.split("\n").length)
                event.data.split("\n").map((data)=> {
                    store.dispatch(quoteActions.updateQuoteSymbolDict(data))
                })
                
                
                // store.dispatch(activityLogAction.appendActivityLog({"id":555}))
                // store.dispatch(wsActions.receiveAllMessages({ messages: event.data }))
            })
            webSocketObj.on("close", () => {})

            // webSocketObj.on(ChatEvent.SendAllMessages, (messages: ChatMessage[]) => {
            //     store.dispatch(chatActions.receiveAllMessages({ messages }))
            // })

            // webSocketObj.on(ChatEvent.ReceiveMessage, (message: ChatMessage) => {
            //     store.dispatch(chatActions.receiveMessage({ message }))
            // })
        }
        next(action)
    }
}
