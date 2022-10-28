import axios from "axios";

axios.defaults.xsrfCookieName = "_gorilla_csrf";
axios.defaults.withCredentials = true;

export default axios;
