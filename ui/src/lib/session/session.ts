import type { components } from "$lib/oapi/gen/types"

export type User = components["schemas"]["User"]

export interface Session {
    initialized: boolean
    authenticated: boolean
    user?: User
}