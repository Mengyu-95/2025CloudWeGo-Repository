namespace go api

struct User {
    1: i64 id
    2: string name
    3: i32 age
    4: string email
}

//struct GetUserRequest {
//    1: i64 user_id
//}

struct GetUserRequest {
    1: i64 id   // 必须包含这个字段
}

struct GetUserResponse {
    1: User user
}


service UserService {
    // 普通Unary方法
    GetUserResponse GetUser(1: GetUserRequest req)
    
    // Kitex流式方法（客户端流/服务端流/双向流）
    // 使用 "streaming.mode" 注解标记流模式
    GetUserResponse StreamUsers(1: GetUserRequest req) (streaming.mode="bidirectional")
    //stream GetUserResponse StreamUsers()
    //stream GetUserResponse StreamUsers()
}