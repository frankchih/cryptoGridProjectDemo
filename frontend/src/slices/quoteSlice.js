import { createSlice } from "@reduxjs/toolkit"

// todo: 以後改活的
const initialState = {
    quoteSymbolList: ["BTCUSDT", "ETHUSDT", "AVAXUSDT"],
    quoteSymbolDict: {"BTCUSDT": {}, "ETHUSDT": {}, "AVAXUSDT": {}},
    quoteDataList: [],
}

export const quoteSlice = createSlice({
    name: "quote",
    initialState,
    reducers: {
        updateQuoteSymbolDict: (state, action) => {
            const dataJson = JSON.parse(action.payload)
            if (dataJson.hasOwnProperty("symbol")) {
                const symbol = dataJson["symbol"]
                const endTime = dataJson["tradeTime"] 
                const price = dataJson["price"] 
                
                // console.log(symbol, price, endTime)
                state.quoteSymbolDict[symbol] = {
                    symbol, price, endTime
                }
                state.quoteDataList = state.quoteSymbolList.map((qs) => {
                    return state.quoteSymbolDict[qs]
                })
            }
            
        },
    },
})
export const quoteActions = quoteSlice.actions

export default quoteSlice.reducer

/*
{"e":"kline","E":1682146337001,"s":"BTCUSDT","k":{"t":1682146336000,"T":1682146336999,"s":"BTCUSDT","i":"1s","f":3090466164,"L":3090466166,"o":"27340.99000000","c":"27340.99000000","h":"27341.00000000","l":"27340.99000000","v":"0.02800000","n":3,"x":true,"q":"765.54780000","V":"0.00800000","Q":"218.72800000"}}

 */
