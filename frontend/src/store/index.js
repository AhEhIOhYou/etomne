import {createStore} from 'vuex';
import {modelsModule} from '@/store/modelsModule';
import {modelModule} from '@/store/modelModule';
import {registrationModule} from '@/store/registrationModule';
import {authorizationModule} from '@/store/authorizationModule';
import {uploadModule} from '@/store/uploadModule';

export default createStore({
  modules: {
    models: modelsModule,
    model: modelModule,
    registration: registrationModule,
    authorization: authorizationModule,
    upload: uploadModule
  },
})
