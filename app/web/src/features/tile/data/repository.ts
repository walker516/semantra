import axios from '@/utils/axiosClient'

export const fetchHealth = async (): Promise<string> => {
  const res = await axios.get('/health')
  return res.data
}
