import {createStore} from 'vuex';
import {modelsModule} from '@/store/modelsModule';
import {registrationModule} from '@/store/registrationModule';
import {authorizationModule} from '@/store/authorizationModule';
import {uploadModule} from '@/store/uploadModule';

export default createStore({
  modules: {
    models: modelsModule,
    registration: registrationModule,
    authorization: authorizationModule,
    upload: uploadModule
  },
})
