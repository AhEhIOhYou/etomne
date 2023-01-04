import axios from "axios";

export const modelsModule = {
    state: () => ({
        models: [],
        isModelsLoading: false,
        page: 1,
        limit: 1,
    }),
    mutations: {
        setModels(state, models) {
            state.models = models;
        },
        setLoading(state, bool) {
            state.isModelsLoading = bool
        },
        setPage(state, page) {
            state.page = page
        },
    },
    actions: {
        async fetchModels({state, commit}) {
            try {
                commit('setLoading', true);
                const response = await axios.get('/api/model', {
                    params: {
                        _page: state.page,
                        _limit: state.limit
                    }
                });
                commit('setModels', response.data)
            } catch (e) {
                console.log(e)
            } finally {
                commit('setLoading', false);
            }
        },
        async loadMoreModels({state, commit}) {
            try {
                commit('setPage', state.page + 1)
                const response = await axios.get('/api/model', {
                    params: {
                        _page: state.page,
                        _limit: state.limit
                    }
                });
                commit('setModels', [...state.models, ...response.data]);
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