#Config: {
  http: {
    listen_port: 8080
  }
  database: {
    host:     !="" // must be specified and non-empty
    user:     !="" // must be specified and non-empty
    password: !="" // must be specified and non-empty
  }
}

config: #Config

config: database: host: "127.0.0.1"
config: database: user: "hello-cue"
config: database: password: "ch@ng3-m3"