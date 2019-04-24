import * as types from '../constants/actionTypes'

const reducer = (state, action) => {
  switch (action.type) {
    case types.WS_MSG_RECEIVED:
      return {...state, currentFrame: action.payload, frames: [...state.frames, action.payload] }
    default:
      return state
  }
}

export default reducer
