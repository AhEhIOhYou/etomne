import {createStore} from 'vuex';
import {modelsModule} from '@/store/modelsModule';
import {registrationModule} from '@/store/registrationModule';
import {authorizationModule} from '@/store/authorizationModule';
import {editModule} from '@/store/editModule';
import {uploadModule} from '@/store/uploadModule';

export default createStore({
  modules: {
    models: modelsModule,
    registration: registrationModule,
    authorization: authorizationModule,
    edit: editModule,
    upload: uploadModule
  },
})
