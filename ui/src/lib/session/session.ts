import type { components } from "$lib/oapi/gen/types"

export type UserInfo = components["schemas"]["UserInfo"]

export interface Session {
    initialized: boolean
    authenticated: boolean
    user?: UserInfo
}