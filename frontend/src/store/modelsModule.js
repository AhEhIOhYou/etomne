import axios from "axios";
import VueCookies from 'vue-cookies';

export const modelsModule = {
    state: () => ({
        models: [],
        userId: null,
        isModelsLoading: false,
        page: 1,
        limit: 1,
    }),
    mutations: {
        setModels(state, models) {
            state.models = models;
        },
        setLoading(state, bool) {
            state.isModelsLoading = bool;
        },
        setUserId(state, userId) {
            state.userId = userId;
        },
        setPage(state, page) { 
            state.page = page;
        },
    },
    actions: {
        async fetchModels({state, commit}) {
            try {
                commit('setLoading', true);
                if (state.userId === null) {
                    const response = await axios.get('/api/model', {
                        params: {
                            _page: state.page,
                            _limit: state.limit
                        }
                    });
                    commit('setModels', response.data)
                } else {
                    const response = await axios.get('/api/model', {
                        params: {
                            _page: state.page,
                            _limit: state.limit,
                            user_id: state.userId
                        }
                    });
                    commit('setModels', response.data)
                }
            } catch (e) {
                console.log(e);
            } finally {
                commit('setLoading', false);
            }
        },
        async loadMoreModels({state, commit}) {
            try {
                commit('setPage', state.page + 1)
                if (state.userId === null) {
                    const response = await axios.get('/api/model', {
                        params: {
                            _page: state.page,
                            _limit: state.limit
                        }
                    });
                    commit('setModels', [...state.models, ...response.data]);
                } else {
                    const response = await axios.get('/api/model', {
                        params: {
                            _page: state.page,
                            _limit: state.limit,
                            user_id: state.userId
                        }
                    });
                    commit('setModels', [...state.models, ...response.data]);
                }
            } catch (e) {
                console.log(e)
            }
        },
        setPagesToOne({state, commit}) {
            commit('setPage', 1);
        }
    },
    namespaced: true
}