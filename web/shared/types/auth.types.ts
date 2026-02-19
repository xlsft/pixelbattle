export type User = {
    auth_date: number
    first_name: string
    hash: string
    id: number
    last_name: string
    photo_url: string
    username: string
}

export type UserModel = {
    uuid: string
    id: number
    name: string
    nickname: string
    picture: string
}