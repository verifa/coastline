
export interface Session {
    initialized: boolean
    authenticated: boolean
    user?: User
}

export interface UserResponse {
    user: User
}

export interface User {
    id: string
    name: string
    email: string
}