import { writable } from 'svelte/store'
import type { User } from './user'

export const user = writable<User>({
    id: "",
    name: "",
    email: "",
})

// class UserStore {
//     constructor(
//         public user: Writable<User> = writable({
//             id: "",
//             name: "",
//             email: "",
//         }),
//         private http = httpStore()
//     ) {
//         this.http.subscribe((resp) => {
//             if (resp.ok && resp.data) {
//                 this.user.set(resp.data.user)
//             }
//         })
//     }

//     public authenticate() {
//         this.http.get("/authenticate")
//     }

//     public

//     public hasError(): boolean {
//         return get(this.http).error !== undefined
//     }
// }

// const userStore = httpStore<UserResponse>()

// export async function authenticate() {
//     await userStore.get("/authenticate")
//     const resp = get(userStore)
//     if (resp.ok && resp.data) {
//         user.set(resp.data.user)
//     } else if (resp.error) {

//     }
// }


// export default function <User>() {
//     const store = writable<Response<Data>>({
//         fetching: true
//     })

//     async function request(method: string, path: string, params?: Record<string, string>, data?: object) {
//         // Clear the store as we are about to make a new request
//         store.update((value) => {
//             value.fetching = true;
//             value.error = undefined;
//             value.data = undefined;
//             return value;
//         });

//         var url = "http://localhost:3000" + "/api/v1" + path
//         const headers = {
//             "Content-type": "application/json"
//         }
//         // If we received query params, add them to the url
//         if (params) {
//             url = url + "?" + new URLSearchParams(params)
//         }
//         const body = data ? JSON.stringify(data) : undefined

//         var storeResponse: Response<Data> = {
//             fetching: false
//         }

//         fetch(url, {
//             method, body, headers,
//             // TODO: this could be same-site when running on the same site
//             credentials: 'include'
//         }).then((response) => {
//             if (!response.ok) {
//                 response.text().then((text) => {
//                     store.update((value) => {
//                         value.fetching = false
//                         value.status = response.status
//                         value.ok = response.ok
//                         value.text = text
//                         return value
//                     })
//                 })
//             } else {
//                 response.json().then((data) => {
//                     store.update((value) => {
//                         value.fetching = false
//                         value.status = response.status
//                         value.ok = response.ok
//                         value.data = data
//                         return value
//                     })
//                 })
//             }
//         })
//             .catch((error) => {
//                 store.update((value) => {
//                     value.fetching = false
//                     value.error = error
//                     return value
//                 })
//             });

//         store.update(() => {
//             return storeResponse
//         })
//     }

//     return {
//         ...store,
//         get: (path: string, params?: Record<string, string>, data?: object) =>
//             request('GET', path, params, data),
//         post: (path: string, params?: Record<string, string>, data?: object) =>
//             request('POST', path, params, data),
//     };
// }