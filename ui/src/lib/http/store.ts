import { writable } from 'svelte/store'

export interface Response<Data> {
    fetching: boolean
    ok?: boolean
    data?: Data
    status?: number
    text?: string
    error?: Error
}

export function createHttpStore<Data>() {
    const store = writable<Response<Data>>({
        fetching: true
    })

    function request(method: string, path: string, params?: Record<string, string>, data?: object) {
        // Clear the store as we are about to make a new request
        store.update((value) => {
            value.fetching = true;
            value.error = undefined;
            value.data = undefined;
            return value;
        });

        let url = "http://localhost:3000" + "/api/v1" + path
        const headers = {
            "Content-type": "application/json"
        }
        // If we received query params, add them to the url
        if (params) {
            url = url + "?" + new URLSearchParams(params)
        }
        const body = data ? JSON.stringify(data) : undefined

        fetch(url, {
            method, body, headers,
            // TODO: this could be same-site when running on the same site
            credentials: 'include'
        }).then((response) => {
            if (!response.ok) {
                response.text().then((text) => {
                    store.update((value) => {
                        value.fetching = false
                        value.status = response.status
                        value.ok = response.ok
                        value.text = text
                        return value
                    })
                })
            } else {
                response.json().then((data) => {
                    store.update((value) => {
                        value.fetching = false
                        value.status = response.status
                        value.ok = response.ok
                        value.data = data
                        return value
                    })
                })
            }
        })
            .catch((error) => {
                store.update((value) => {
                    value.fetching = false
                    value.error = error
                    return value
                })
            });
    }

    return {
        ...store,
        get: (path: string, params?: Record<string, string>, data?: object) =>
            request('GET', path, params, data),
        post: (path: string, params?: Record<string, string>, data?: object) =>
            request('POST', path, params, data),
    };
}