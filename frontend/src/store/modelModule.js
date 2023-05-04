// import axios from "axios";

// export const modelModule = {
//     state: () => ({
//         model: {},
//         isModelLoading: false,
//     }),
//     mutations: {
//         setModel(state, model) {
//             state.model = model;
//         },
//         setLoading(state, bool) {
//           state.isModelLoading = bool
//       },
//     },
//     actions: {
//         async fetchModel({state, commit}, id) {
//             try {
//                 commit('setLoading', true);
//                 const response = await axios.get(`/api/model/${id}`);
//                 commit('setModel', response.data);
//             } catch (e) {
//                 console.log(e)
//             } finally {
//               commit('setLoading', false);
//           }
//         },
//     },
//     namespaced: true
// }