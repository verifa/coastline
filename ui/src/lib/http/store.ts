import { writable } from 'svelte/store'

export interface Response<Data> {
    fetching: boolean
    ok: boolean
    // requested is marked as true after at least one request has been made
    requested: boolean
    data?: Data
    status?: number
    text?: string
    error?: Error
    redirected?: boolean
}

export function createHttpStore<Data>() {
    const store = writable<Response<Data>>({
        fetching: false,
        ok: false,
        requested: false,
    })

    function request(method: string, path: string, params?: Record<string, string>, data?: object) {
        // Clear the store as we are about to make a new request
        store.update((value) => {
            value.ok = false
            value.fetching = true;
            value.requested = true;
            value.error = undefined;
            value.data = undefined;
            return value;
        });

        let url = "/api/v1" + path
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
            // TODO: this could be same-origin when running on the same site
            credentials: 'include'
        }).then((response) => {
            if (response.ok) {
                if (response.redirected) {
                    store.update((value) => {
                        value.ok = response.ok
                        value.status = response.status
                        value.redirected = response.redirected
                        return value
                    })
                    return
                }
                response.json().then((data) => {
                    store.update((value) => {
                        value.fetching = false
                        value.status = response.status
                        value.ok = response.ok
                        value.data = data
                        return value
                    })
                }).catch((error) => {
                    store.update((value) => {
                        value.fetching = false
                        value.error = error
                        value.text = "Error converting response to JSON"
                        return value
                    })
                })
            } else {
                response.text().then((text) => {
                    store.update((value) => {
                        value.fetching = false
                        value.status = response.status
                        value.ok = response.ok
                        value.text = text
                        return value
                    })
                }).catch((error) => {
                    store.update((value) => {
                        value.fetching = false
                        value.error = error
                        value.text = "Error converting response to text"
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