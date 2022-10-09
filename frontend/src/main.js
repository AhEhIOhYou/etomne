import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
// import UI from '@/UI';
import components from '@/components';
import directives from '@/directives';
const app = createApp(App);

app.component(components.ModelsItem.name, components.ModelsItem);
app.component(components.ModelsList.name, components.ModelsList);
app.component(components.Navbar.name, components.Navbar);
app.directive(directives.VIntersection.name, directives.VIntersection);

app.use(store).use(router).mount('#app');
