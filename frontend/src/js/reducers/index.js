import * as types from '../constants/actionTypes'

const reducer = (state, action) => {
  switch (action.type) {
    case types.WS_MSG_RECEIVED:
      state.game.move(action.payload.history[action.payload.history.length-1])
      return {...state, ...action.payload }
    case types.MK_MOVE:
      return {...state, history: action.payload }
    case types.CH_NAME:
      return {...state, username: action.payload }
    case types.START_GAME_PENDING:
      return {...state, searching: true }
    case types.START_GAME_FULFILLED:
      return {...state, ...action.payload }
    default:
      return state
  }
}

export default reducer
