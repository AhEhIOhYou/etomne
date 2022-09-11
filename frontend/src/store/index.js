import {createStore} from 'vuex';
import {modelsModule} from '@/store/modelsModule';
import {modelModule} from '@/store/modelModule';

export default createStore({
  modules: {
    models: modelsModule,
    model: modelModule,
  },
})
