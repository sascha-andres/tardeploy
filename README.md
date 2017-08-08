# tardeploy

A tool to automate deployment of web applications onto a web server

## Configuration

Configuration is split into three main parts:

1. Directories
2. Application
3. Trigger

Directories are the directories being worked in while application controls some aspects of the behaviour. Trigger are programs that are executed before and after a deployment.

The configuration is done using a YAML, TOML or ini file. This readme shows all samples as YAML.

The configuration file can be located in one of three locations:

1. /etc/tardeploy/tardeploy.yaml
2. $HOME/.tardeploy/tardeploy.yaml
3. ./tardeploy.yaml

A sample configuration looks like this:

    ---
    
    Directories:
      TarballDirectory: /home/andres/tardeploy/tarballs/
      WebRootDirectory: /home/andres/tardeploy/webroot
      ApplicationDirectory: /home/andres/tardeploy/apps
      Security:
        User: andres
        Group: andres
    Application:
      NumberOfBackups: 1
      BatchInterval: 10
      LogLevel: info
    Trigger:
      Before: /bin/echo

### Directories

#### TarballDirectory

This is the directory where the daemon looks for tarballs to deploy. A tarball must have an extension of `.tgz` or `.tar.gz`.

#### WebRootDirectory

A directory where the links to the deployed applications will be created. For a tarball named `app.tgz` a symbolic link `app` will be created.

#### ApplicationDirectory

This is the place where the extracted applications reside. Each directory has it's own directory containing 1 to n folders with different versions.

#### Security

After an application is extracted, the file owner has to be changed to the webserver user.

##### User

The web server's user account

##### Group

The web server's group

### Application

#### NumberOfBackups

The number of backups to keep in the application directory. Set to -1 to have all versions retained, 0 to have no backups or any higher number to retain a certain amount of versions.

#### BatchInterval

The daemon batches changes to tarballs as a protection against constantly changing files. You can change the amount of time in seconds to push changed tarballs to the runtime.

#### LogLevel

You can set the log level to one of the four log levels:

1. debug
2. info
3. warn
4. error

Everything higher or equal to the configured level will be logged. Eg if you set it to `warn`, `error` will be logged too, but not `debug` or `info`.

#### TarCommand

If you do not want to use the internal tar.gz handling, set this option to the path of the system tar executable. The tarballs will then be executed with he following commandline:

    <TarCommand> xzf <tarball>

The command is run in the deployment directory ( that is the subfolder within the application directory which might be kept as a backup ).

### Trigger

Triggers are a way to interact with other systems. A trigger will be executed as a subprocess. A trigger will be called with two parameters:

1. application
2. trigger name

#### Before

The before trigger will be executed before any work will be done

#### After

The after trigger will be executed after the deployment has been done

## Hints

You might want to consider that the folder within the webroot is a symbolic link. Some webserver software requires special options to be set to support that type of configuration.

## systemd installation

A sample service file can be found in the directory `systemd/`. It assumes you have the binary of tardeploy in `/opt/tardeploy/bin/tardeploy`

## History

|Version|Description|
|---|---|
|0.1.0|Initial version|