# func

Turn "stuff" into some kind of standard binary (a cli) that can act as a "lambda". This is done by generating up to 2 cobra commands for each provided "stuff". It's up to the one providing the "stuff" to create the actual binary (the cli). (See examples)

*From my non-empirical studies on the example lambdas in the test folder, I found out that most of the time during execution is spent outside the provided function. Ie in the fluff, (de)serializing bytes, while fast, still takes time.*

PS. but I've heard, that the main benefit of "lambdas" are the dx.

## Generated CLI

The generated cli would have a couple of commands. `call`-command that would expect you to pass the request into stdin and would return the result to stdout.

`listen`-command that would start a server that either listens to a .sock-file or a tcp port.

### Carry around an extra command

2 commands in one binary might seem like a waste of a good binary. But they are pretty small anyway.

## Helpers

There's a number of helpers created. Helpers to wrap functions and wendy modules and turn them into servers etc. But also helpers to call the result of the binaries. In form of a very generic client, and something that can execute commands, feed them stuff on stdin and capture their output.

## Wendy

Turning wendy handlers and modules into "lambdas". This turned out to be a pretty good fit.

Some infrastructure to route messages are still needed, but idealy you'd shove all messages on the same queue and add workers as needed.

### TODO

* I think the decoder_for_loop does what's expected, but it would be nice WITH SOME UNIT TESTS...
* Logging should be captured both from the framework and from the "lambda"... somewhere, somehow.
* Some stats about each call should be recorded somewhere, somehow.

## HTTP

... http "lambdas". Turning stuff from http package into a "lambda".

Turning a single http.HandlerFunc into a callable lambda works pretty well. And it can be wrapped for days to add functionality. Unfortunately, parts of the http api leaves a lot to wish for. And so does http servers to some extent. (See HTTP servers below)

### TODO

* Logging should be captured both from the framework and from the "lambda"... somewhere, somehow.
* Some stats about each call should be recorded somewhere, somehow.

### HTTP Servers

Both http.Serve and gin.RunX has a slight start up penalty. Using them with the call command gives you that penalty every call. It is prolly better to teach helper/starters.Gin() to also serve on an unix socket to "lambdaify" them that way.

## GRPC

...

I expect grpc servers to have the same penalty as http servers did. That protocol they prefer to run on thrives surrounded by tls if I've understood it correctly. Then on the other hand, they're usually only interfaces. Sky is the limit.

I've foolishly been avoiding GRPC so far, but they're on the "roadmap".

## Other

...

### RPC?

> While the api is simple, it has twisted a few knobs too much towards nats. There are things in that context that will be very hard to implement in a `call`. Perhaps turned into a service somehow it could be done. But still sketchy, and wasted.

### EWF?

> This old scala project of mine suddenly show some promise. However, wendy (which is made in Go, duh(!)) is close to the same thing, and at least as competent, in particular in the response department. Then again, perhaps there's still room for a go-ewf over wendy? ^^

### Otto?

> This is intresting territory. As much as I hate javascript, it is still kind of useful at times.

Like hosting sveltekit apps, build a custom adapter and then run them on otto. Wonder if that would fly... fast? In worst case, im sure there's node.js bindings somewhere. But then you're prolly in cgo territory. Aaand if you've invested in running the custom binaries created by this lib, then what's the really the cost of running a different binary?