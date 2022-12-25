import axios from "axios";

const auth = () => {
  axios.post('/api/users/refresh', {
    refresh_token: refreshToken
  })
  .then(response => {
    $cookies.set('access_token', response.data.access_token, '15m', '/');
    $cookies.set('refresh_token', response.data.refresh_token, '7d', '/');
    return { accessToken: response.data.access_token, refreshToken: response.data.refresh_token }
  })
  .catch(error => {
    console.log(error);
  })
};

export default auth;

// TODO:
// Засунуть в переменную и использовать во всех запросах