#  Actions

When a deployment is executed against a host(s) typically one or more **actions** will be performed against that host in order to configure as expected. This document details the **built-in** actions, however to extend the functionality of [plunder](github.com/plunder-app/plunder) there is the capability to extend the available actions through the use of plugins. 



## Built-in Actions

All actions are defined by a `type` which specifies what tasks the action will perform, also all actions should come with a `name` that identifies what the action will perform. The names should make it easy to identify relevant tasks as they're executed or when selecting individual tasks when using the user Interface.

Example in json and yaml below:

```json
{
  "task" : "command",
  "command" : "docker run image",
  "name" : "Starts the docker image \"image\""
}
```

```yaml
- task: download
  source: "/home/user/my_archive.tar.gz"
  name: "Retrieve the home archive"
```

### Command

The **command** action type is used to execute a command either locally or remote, it will exit execution if the command fails (or it can be ignored) and the results can be stored to be executed at a later point. 

Set the `ignoreFail` to `true` to allow execution of tasks to continue in the event that the command fails. If a long running task should be known to only execute for a specific amount of time, commands can be given a timeout which will end the command should it not complete in time. The `timeout` setting should be set in seconds which will specify how long the task is allowed to execute for.

```yaml
- task: command
  command: "sleep 100"
  timeout: 99
  ignoreFail: true
```

 *The above example will execute a sleep for a hundred seconds, however the command has a timeout set for only 99 seconds. Execution will be halted once the timeout is met, and if the task returns a fail code the execution will continue onto the next action*

If a command requires elevated privileges, the `commandSudo` option allows executing a command as different user, with it's entitled privileges. 

**Note**: This requires `NOPASSWD` to be set for the current user.

```yaml
{
  "task" : "command",
  "command" : "cat /dev/null > /var/log/messages",
  "name" : "Concatenate the messages file to clear space",
  "commandSudo" : "root"
}
```

#### Using commands between actions deployments

There may be a requirment to save the output of a command to be used in a different action or a different deployment, some commands will generate tokens or output that can be used at a later point.

There are two options to save the output of a command: 

- `commandSaveFile` - saves the command output to a path
- `commandSaveAsKey` - Saves the ouput in-memory under a specified `key`

These saved ouputs can then be used later through the use of the `key` options:

- `KeyFile` - executes the commands in the file specified under the `path`

- `KeyName` - executes the commands saved in-memory under the specified `key`

  

The below example will create a command Key under the name `joinKey` (JSON format) :

```json
{
	"name" : "Generate a join token",
  "type" : "command",
  "command" : "kubeadm token create --print-join-command 2>/dev/null",
  "commandSaveAsKey" : "joinKey" 
}
```



This key can now be used in a different deployment with different hosts (YAML format):

```yaml
- type: "command"
  name: "Join worker to Kubernetes cluster"
  keyName: "joinKey"
  commandSudo : "root"
```



#### Piping data between commands

In the event that data needs to piped into a remote command the options `commandPipeFile` and `commandPipeCmd` can be used. The first will take the contents of `path` and pass it as `STDIN` to the command being executed under the option `command`. The `commandPipeCmd` will execute a command locally and pass the `STDOUT` of that command into the `STDIN` of the command being ran under the `command` option.



The below example will run the command `echo "deb https://apt.kubernetes.io/ kubernetes-xenial main"` locally, and pass the `STDOUT` to the command `tee /etc/apt/sources.list.d/kubernetes.list` that is being ran using `sudo` privileges.

```yaml
  - type: command
    command: "tee /etc/apt/sources.list.d/kubernetes.list"
    commandPipeCmd: echo "deb https://apt.kubernetes.io/ kubernetes-xenial main"
    name: Set Kubernetes Repository
    commandSudo: root

```

This is useful for a variety of usecases, although it has been found very useful for appending data to existing files that require elevated privilieges. 

**Example reasons for piping data to a command**

The command `echo "SOME data" | tee /a/file/that/needs/sudo/privs` will fail even with `commandSudo`, the reason for this is that the `sudo` is only going to work for everything upto the pipe. The remaining part of the command will be ran as the current user and therefore doesn't have the required privileges.

### Upload / Download of files

Both of the command types `upload` and `download` have the same set of options:

- `destination` - Where the file will be once the `upload`/`download` has completed
- `source` - The file that will be either `uploaded`/`downloaded`
- `name` - Details what the action will be doing

```yaml
  - type: download
    destination: ./ubuntu.tar.gz
    name: Retrieve local copy of ubuntu.tar.gz
    source: ./ubuntu.tar.gz
```



### Plugins

Plugins allow the creation of unique actions to be performed, such as specific interactions with platforms, programs and infrastructure. All plugins will load at startup and register their actions into the parlay engine. Passing information to a plugin should be done in the following manner:

```yaml
  - name: Push kubernetes images for managers
    plugin:
      imageName:
      - k8s.gcr.io/kube-apiserver:v1.14.0
      - k8s.gcr.io/kube-controller-manager:v1.14.0
      - k8s.gcr.io/kube-scheduler:v1.14.0
      localSudo: true
      remoteSudo: true
    type: docker/image
```

The main differences are:

- `plugin` - Contains all of the specifics that will be passed to the plugin logic
- `type` - Should be the action defined by the plugin itself.





## Additonal configuration

### No Password sudo

To enable password-less sudo the `/etc/sudoers` file needs modifying (DO NOT DO THIS MANUALLY).

To edit the sudo file use the following command:

```
sudo visudo
```

Then add the following entry to the end of the file, replacing the `username` with the correct entry :

```
username     ALL=(ALL) NOPASSWD:ALL
```

This can be tested by either opening a new session or logging out and back in and then testing that `sudo <CMD>` doesn't require a password.
