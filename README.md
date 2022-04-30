# go-gurl

![Demo](assets/demo.gif)

TUI `"curl"` as a personal replacement of Postman. Postman just feels overly bloated in appearance and as a program.
This is written purely in Golang and does not use `curl` under the hood.

> I'm working on this while learning Go.

⚠️ This program has lots of problems. Most I may not even know about. But the following are the known limitations at the
moment:

- [ ] Windows support (I don't have access to a windows machine at the moment)
- [ ] Cookies support
- [ ] Headers deletion
- [ ] Save/Manipulate response/request
- [ ] Edit/delete bookmarks and history
- [ ] Tests (I have some written)
- [ ] Commands
- [ ] Input validations

I'll be providing implementations as time goes by. Feel free to submit an issue as you try this out.

## How to use

### Installation

```bash
go install github.com/blackmann/go-gurl
```

### Keybinds

| Bind        | Action                                                                                                                    |
|-------------|---------------------------------------------------------------------------------------------------------------------------|
| `shift+tab` | Alternate between views (address bar and viewport)                                                                        |
| `esc`       | Enter/leave command mode. <br/>In command mode, you can press the forward or back key to switch between the viewport tabs |
| `ctrl+c`    | Quit                                                                                                                      |
| `$`         | Show history. You can filter history with ID or annotation. <br/>_See below on how to annotate history_                   |

### History

When requests are made, they are saved into history. To trigger the history modal, enter a leading `$`. 
You can filter the history with ID number or annotation. Selecting a history item will prefill all request fields (address, headers and body).

To annotate history, first find the history ID then enter command mode (`esc`) then type

```
/annotate $32 create-account
```

This feature is useful when you run a request very often when testing.

### Bookmarks

Bookmarks allow to you create and use alias for base paths/endpoints. 
For example, if you mostly work with an endpoint `https://jsonplaceholder.typicode.com`, you can be able to create a bookmark with (in command mode `esc`)

```
@typicode https://jsonplaceholder.typicode.com
```

You can then use the bookmark in making requests (in the addressbar as)

```
POST @typicode/todos/
```


