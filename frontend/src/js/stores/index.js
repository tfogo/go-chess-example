import { createStore, applyMiddleware } from 'redux'
import { createLogger } from 'redux-logger'
import thunk from 'redux-thunk'
import reducers from '../reducers'
import { initWebSocket } from '../actions'

const defaultState = {
  currentFrame: {},
  frames: []
}

const middleware = applyMiddleware(createLogger(), thunk)
const store = createStore(reducers, defaultState, middleware)
initWebSocket(store.dispatch)
//store.dispatch(initWebSocket)
window.store = store // expose store globally to manipulate in browser

export default store
