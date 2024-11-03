package dtos

type ClerkCreateUserDto struct {
    Data struct {
        ID string `json:"id"`
    } `json:"data"`
}


type UserResponseDto struct {
    Username string `json:"username"`
    Email    string `json:"email"`
}