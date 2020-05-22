# AppPauser

[![forthebadge](https://forthebadge.com/images/badges/oooo-kill-em.svg)](https://forthebadge.com)

*an application pauser for all occasions*

This application was inspired by [this post on Reddit](https://www.reddit.com/r/pcgaming/comments/e7ftzr/unpausable_cutscenes_i_made_a_windows_application/), which showcased an application that allows one to pause a game during unpausable cutscenes on Windows. 
OP's comments showed that the application had to pause the executing thread to pause a game, which led me to thinking about developing a similar application for Linux. 
It's indeed very simple to send STOP and CONT signals on Linux, which made it possible to create a launcher wrapper that listens on a socket for commands like `pause`, `resume` and `kill`.

## How to use

Launch the desired application with 

```shell
$ apppauser your-application [arguments]
```

and control it with

```shell
apppauserctl command
```

where `command` can be:

* `pause`: sends the SIGSTOP signal to an application
* `resume`: sends the SIGCONT signal to an application
* `toggle`: sends either SIGSTOP or SIGCONT signal to an application, depending on its current status
* `kill`: sends the SIGTERM signal to an application

## How to install

Install the Go compiler, and then run `go get github.com/ericonr/AppPauser` to compile both binaries and store them in `$GOBIN`. 
You then only need to add `$GOBIN` to your `$PATH` and you are ready to use this!

## Future steps

There's still a lot to do. What I've already thought about:

* Make the default name for sockets be related to the user's username
* Include feedback on the output of apppauserctl, which for now doesn't listen to any responses
* Create command line options to determine the socket name and other options (such as starting an application already paused)

## Warning

This is as incomplete as it comes, but it did work in my machine. Please let me know about any issues!
