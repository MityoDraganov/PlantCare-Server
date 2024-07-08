func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Hash the password (simplified, use a proper library in production)
    user.PasswordHash = hashPassword(user.PasswordHash)

    users = append(users, user)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func hashPassword(password string) string {
    // This is a placeholder. In production, use a library like bcrypt.
    return fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for _, user := range users {
        if user.Username == params["username"] {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(user)
            return
        }
    }
    http.Error(w, "User not found", http.StatusNotFound)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, user := range users {
        if user.Username == params["username"] {
            users = append(users[:index], users[index+1:]...)
            var updatedUser User
            err := json.NewDecoder(r.Body).Decode(&updatedUser)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }

            // Hash the password (simplified, use a proper library in production)
            updatedUser.PasswordHash = hashPassword(updatedUser.PasswordHash)
            users = append(users, updatedUser)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedUser)
            return
        }
    }
    http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    for index, user := range users {
        if user.Username == params["username"] {
            users = append(users[:index], users[index+1:]...)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(users)
            return
        }
    }
    http.Error(w, "User not found", http.StatusNotFound)
}
