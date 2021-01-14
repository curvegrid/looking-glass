import axios from 'axios';
import { baseUrl } from './constants';

export default axios.create({
  withCredentials: true,
  baseURL: baseUrl,
});
