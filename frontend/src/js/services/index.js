import axios from 'axios'

axios.defaults.headers.post['Content-Type'] = 'application/json'
axios.defaults.baseURL = process.env.CONTACTS_API_URI


const Api = {
  fetchVideos: () => {
    return axios({
      method: 'get',
      url: `/videos`
    })
  }
}

export default Api
