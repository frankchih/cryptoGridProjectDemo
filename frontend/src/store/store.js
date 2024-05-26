import { applyMiddleware, combineReducers, configureStore } from "@reduxjs/toolkit"
import activityLogReducer from "../slices/activityLogSlice"
import { webSocketMiddleware } from "../middlewares/ws"
import wsSlice from "../slices/wsSlice"
import quoteSlice from "../slices/quoteSlice"

import { getLeverageAppApi } from "../services/leverageAppRTKService"
import { setupListeners } from "@reduxjs/toolkit/dist/query"
import { getMainAppApi } from "../services/mainAppRTKService"
import { getOrderAppApi } from "../services/orderAppRTKService"


export const store = configureStore({
    reducer: {
        activityLogReducer: activityLogReducer,
        wsSlice: wsSlice,
        quoteSlice: quoteSlice,
        [getLeverageAppApi.reducerPath]: getLeverageAppApi.reducer, 
        [getMainAppApi.reducerPath]: getMainAppApi.reducer,
        [getOrderAppApi.reducerPath]: getOrderAppApi.reducer, 
    },
    middleware: (getDefaultMiddleware) => {
        return getDefaultMiddleware().concat([
            getOrderAppApi.middleware,
            getLeverageAppApi.middleware,
            getMainAppApi.middleware,
            webSocketMiddleware,            
        ])
    },
})

setupListeners(store.dispatch)
