import { createStore, applyMiddleware } from 'redux'
import { createLogger } from 'redux-logger'
import thunk from 'redux-thunk'
import reducers from '../reducers'
import { initWebSocket } from '../actions'
import Chess from "chess.js";

const defaultState = {
  username: "",
  white: "",
  black: "",
  gameId: "",
  gameStarted: false,
  searching: false,
  game: new Chess(),
  history: []
}

const middleware = applyMiddleware(createLogger(), thunk)
const store = createStore(reducers, defaultState, middleware)
//initWebSocket(store.dispatch)
//store.dispatch(initWebSocket)
window.store = store // expose store globally to manipulate in browser

export default store
