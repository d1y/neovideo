import axios from 'axios'

const srv = axios.create({
  baseURL: '/api/v1',
})

export default srv
