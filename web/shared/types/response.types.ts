export type ResponseErrorConstructor = {
    cause?: any;
    message: string;
    code: number;
}

export type ResponseError = {
    error: ResponseErrorConstructor
} | undefined