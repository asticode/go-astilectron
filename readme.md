# Intro

TODO

# Tutorials

Chances are what you'd like to implement has already been implemented in one of [the examples](https://github.com/asticode/go-astilectron/tree/master/examples).

To run any of the example run the following command:

    $ go run examples/<name of the example>/main.go -v
    
Here's a list of the examples:

- **1.basic_window** creates a basic window that displays a static .html file
- **2.basic_window_events** plays with basic window methods and shows you how to set up your own listeners
- **3.webserver_app** shows you how simple it is to set up a webserver app with **astilectron**

# Features and roadmap

- [x] window basic methods (create, show, close, resize, minimize, maximize, ...)
- [x] window basic events (close, blur, focud, unresponsive, crashed, ...)
- [ ] remote messaging (messages between GO and the JS in the webserver)
- [ ] menu methods
- [ ] menu events
- [ ] session methods
- [ ] session events
- [ ] window advanced options (add missing ones)
- [ ] window advanced methods (add missing ones)
- [ ] window advanced events (add missing ones)
- [ ] child windows