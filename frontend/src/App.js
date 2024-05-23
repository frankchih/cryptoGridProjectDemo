import React from "react"
import { Provider } from "react-redux"

import Main from "./components/Main"
import "./index.css"
import { store } from "./store/store"

const App = () => {
    return (
        // <React.StrictMode>
            <Provider store={store}>
                <Main />
            </Provider>
        // </React.StrictMode>
    )
}

export default App
