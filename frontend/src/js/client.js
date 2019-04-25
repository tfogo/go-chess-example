import React from "react"
import ReactDOM from "react-dom"
import App from "./components/App.jsx"
import store from './stores'
import { Provider } from 'react-redux'
import '../index.html'
import '../styles/main.css'

ReactDOM.render(
  <Provider store={store}>
    <div className="h-100">
      <App></App>
      </div>
  </Provider>,
  document.getElementById('app')
)
