import { createApp } from 'vue'
import App from './App.vue'
import './assets/css/tailwind.css'
import './assets/css/theme.css'
import { FontAwesomeIcon } from './plugins/fontawesome'
// Для Font Awesome
declare global {
  interface Window {
    FontAwesomeConfig: {
      autoReplaceSvg: string;
    };
  }
}

window.FontAwesomeConfig = {
  autoReplaceSvg: 'nest'
};
const app = createApp(App);

// Регистрация глобального компонента Font Awesome
app.component('font-awesome-icon', FontAwesomeIcon);

app.mount('#app');