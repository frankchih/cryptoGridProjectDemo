import { createSlice } from "@reduxjs/toolkit"


const initialState = {
    messages: [],
    isEstablishingConnection: false,
    isConnected: false,
}

export const wsSlice = createSlice({
    name: "ws",
    initialState,

    reducers: {
        startConnecting: (state) => {
            console.log("slice in startConnecting")
            state.isEstablishingConnection = true
        },
        connectionEstablished: (state) => {
            console.log("slice in connectionEstablished")
            state.isConnected = true
            state.isEstablishingConnection = true
        },
        receiveAllMessages: (state, action) => {
            state.messages = action.payload.messages
        },
        receiveMessage: (state, action) => {
            state.messages.push(action.payload.message)
        },
        submitMessage: (state, action) => {
            return
        },
    },
})

export const wsActions = wsSlice.actions

export default wsSlice.reducer
