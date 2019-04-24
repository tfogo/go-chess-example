import React from "react"
import ReactDOM from "react-dom"
import App from "./components/App.jsx"
import store from './stores'
import { Provider } from 'react-redux'
import { BrowserRouter as Router, Route, Link } from 'react-router-dom'
import '../index.html'
import '../styles/main.css'

ReactDOM.render(
  <Provider store={store}>
    <Router>
    <div className="h-100">
      <Route exact path="/" component={App} />
      </div>
    </Router>
  </Provider>,
  document.getElementById('app')
)
