import type { components } from "$lib/oapi/spec"

export type UserInfo = components["schemas"]["UserInfo"]

export interface Session {
    initialized: boolean
    authenticated: boolean
    user?: UserInfo
}