## What is hot update?

If you have used nginx, you may know that nginx supports hot update, which means you can update your nginx without stopping and restarting it. It serves old connections with old version, and accepts new connections with new version. Notice that hot compiling is different from hot update, where hot compiling is monitoring your source files and recompile them when the content changes, it requires stop and restart your applications, `bee start` is a tool for hot compiling.


## Is hot update necessary?

Some people says that hot update is not as useful as its cool name. In my opinion, this is absolutely necessary because zero-down server is our goal for our services. Even though sometimes some errors or hardware problems may occur, but it belongs to design of high availability, don't mix them up. Service update is a known issue, so we need to fix this problem.


## How Beego support hot update?

The basic principle of hot update: main process fork a process, and child process execute corresponding programs. So what happens? We know that after forked a process, main process will have all handles, data and stack, etc, but all handles are saved in `CloseOnExec`, so all copied handles will be closed when you execute it unless you clarify this, and we need child process to reuse the handle of `net.Listener`. Once a process calls exec functions, it is "dead", system replaces it with new code. The only thing it left is the process ID, which is the same number but it is a new program after executed.

Therefore, the first thing we need to do is that let child process fork main process and through `os.StartProcess` to append files that contains handle that is going to be inherited.

The second step is that we hope child process can start listening from same handle, so we can use `net.FileListener` to achieve this goal. Here we also need FD of this file, so we need to add a environment variable to set this FD before we start child process.

The final step is that we want to serve old connections with old version of application, and serve new connections with new version. So how can we know if there is any old connections? To do this, we have to record all connections, then we are able to know. Another problem is that how to let new application to accept connection? Because two versions of applications are listening to same port, they get connection request randomly, so we just close old accept function, and it will get error in `l.Accept`.

Above are three problems that we need to solve, you can see my code logic for specific implementation.


## Show time

1. Write code in your Get method:

		func (this *MainController) Get() {
			a, _ := this.GetInt("sleep")
			time.Sleep(time.Duration(a) * time.Second)
			this.Ctx.WriteString("ospid:" + strconv.Itoa(os.Getpid()))
		}

2. Open two terminals:

	One execute: ` ps -ef|grep <application name>`

	Another one execute： `curl "http://127.0.0.1:8080/?sleep=20"`

3. Hot update

	`kill -HUP <PID>`

4. Open a terminal to request connection: `curl "http://127.0.0.1:8080/?sleep=0"`

As you will see, the first request will wait for 20 seconds, but it's served by old process; after hot update, the first request will print old process ID, but the second request will print new process ID.
