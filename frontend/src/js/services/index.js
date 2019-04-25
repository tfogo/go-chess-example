import axios from 'axios'

axios.defaults.headers.post['Content-Type'] = 'application/json'
axios.defaults.baseURL = "http://localhost:9001"


const Api = {
  startGame: (username) => {
    return axios({
      method: 'post',
      url: `/start`,
      data: {
        username
      }
    })
  }
}

export default Api
