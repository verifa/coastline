import { writable } from 'svelte/store'
import type { Session, UserInfo } from './session'

export const session = createSessionStore()

function createSessionStore() {
    const { subscribe, set } = writable<Session>({
        initialized: false,
        authenticated: false
    });

    return {
        subscribe,
        login: (user: UserInfo) => {
            set({
                initialized: true,
                authenticated: true,
                user: user
            })
        },
        logout: () => {
            set({
                initialized: true,
                authenticated: false,
                user: undefined
            })
        },
    };
}