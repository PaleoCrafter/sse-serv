# sse-serv - 0.0.0

```
                                                               ┌───────────┐
 ┌─────────────┐   ┌────────────┐   ┌───────────┐ <─ register ─┤  client   │
 │             │<──┤            │<──┤           ├──── SSE ────>│ [browser] │   
 │ AMQP broker │   │  sse-serv  │   │ rev proxy │              └───────────┘ 
 │ [RabbitMQ]  │   │   [this]   │   │  [nginx]  │              ┌───────────┐
 │             ├──>│            ├──>│           ├──── SSE ────>│  client   │ 
 └─────────────┘   └────────────┘   └───────────┘ <─ register ─┤ [browser] │        
                                                               └───────────┘
                                                                    ...
```

*Do not use in production.*

## License
[MIT](https://opensource.org/licenses/MIT) - © dtg [at] lengo [dot] org
