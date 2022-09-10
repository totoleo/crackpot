namespace go a.b.c

struct Request {
    1:string Name
}

struct Message {
    1:string message
}
service Echo {
    Message Hello(1:Request req)
}