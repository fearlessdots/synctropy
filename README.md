# Synctropy

```
 ____                   _
/ ___| _   _ _ __   ___| |_ _ __ ___  _ __  _   _
\___ \| | | | '_ \ / __| __| '__/ _ \| '_ \| | | |
 ___) | |_| | | | | (__| |_| | | (_) | |_) | |_| |
|____/ \__, |_| |_|\___|\__|_|  \___/| .__/ \__, |
       |___/                         |_|    |___/
```

`synctropy` is a user-friendly wrapper that simplifies the synchronization and management of configurations using tools like unison and rsync. With `synctropy`, you can easily create and manage customized crates, which act as containers for organizing and synchronizing specific sets of configurations. Within each crate, you can define targets, representing the individual configurations you want to synchronize.

At the heart of `synctropy` lies its modular approach, enabled by hooks. Hooks are an essential component that allows you to run custom scripts or commands at specific stages of various operations within `synctropy`. This powerful feature empowers you to implement your custom logic to the program. Hooks provide the flexibility to extend and customize the behavior of not only the synchronization process but also other operations such as viewing and editing configurations, making `synctropy` a truly modular program.

With `synctropy`, you have the freedom to tailor not only the synchronization process but also other aspects of the program to suit your requirements. By leveraging hooks, you can incorporate additional functionality, perform complex transformations, or integrate with external systems. The modular nature of `synctropy` ensures that it can adapt to various use cases and provide a flexible solution for managing and synchronizing configurations, as well as other operations within the program.

Bearing a name that fuses `sync` and the scientific concept `syntropy` - signifying the shift from disorder to structure, `synctropy` aims to manage the mix of your various files and turn them into a smoothly synchronized collection. It's about evolving from entropy to syntropy, converting the disordered into the organized.

## Key Features

- **Crates and Targets:** `synctropy` provides a flexible and organized structure for managing data synchronization. Crates represent a collection of targets that can be used to control the syncing process via, for example, the `pre_transaction` and `post_transaction` hooks. It can also hold configuration and other files used across most or all of its targets. The targets also provide their own set of hooks and are the structural elements actually used to sync files and directories.

- **Hooks Support:** Both crates and targets support hooks, allowing you to run custom scripts or commands before or after syncing. The sync process itself is also defined by a hook. This enables you to perform additional actions or customize the synchronization process according to your specific needs. These hooks provide a modular way to extend and customize the synchronization process according to your specific needs. With hooks, you can seamlessly integrate additional functionality, perform complex transformations, or interact with external systems.

- **Templates Support:** Creating new crates and targets is made easier with template support. Templates provide a convenient way to create consistent configurations by predefining common settings and hooks.

- **Utilities for Hooks:** `synctropy` includes a set of utilities that can be used with hooks to simplify common tasks. These utilities allow you to display messages, show section titles, print colored output, ask for user confirmation, and more. They enhance the functionality of hooks and provide a convenient way to interact with the user during the syncing process.

- **Command-Line Interface and User-Friendly Structure:** With its command-line interface, `synctropy` offers a straightforward and efficient way to manage your synchronization tasks. The program is designed with a user-friendly structure, providing intuitive commands and comprehensive documentation.

## Use Cases

Synctropy can be used in a wide range of scenarios to simplify and automate synchronization tasks. Here are some use cases where Synctropy can be beneficial:

### Personal File Syncing

Synctropy can be used in conjunction with programs like rsync and Unison to sync personal files between different devices or cloud storage services. By defining `target` configurations that specify the source and destination paths, you can easily keep your files in sync, ensuring that you have the latest versions available on all your devices.

### Backup and Restore

Synctropy can be used as a backup and restore tool for important files or directories on remote systems. By creating `target` configurations that specify the source files or directories, you can easily back up the data to another location. In case of data loss or system failure, you can then use Synctropy to restore the backed-up files, ensuring the availability and integrity of your data.

### Distributed Systems Management

Synctropy could also be useful for managing and synchronizing configurations across multiple distributed systems. Whether you have a cluster of servers, a network of IoT devices, or a fleet of containers, Synctropy can help ensure consistency and reduce manual effort by synchronizing configurations across all the systems. This helps maintain a unified state and simplifies the management of distributed environments.

## Installation

> I have made a PKGBUILD available in this repository, which allows for easy building and installation on Arch Linux.

To build the program, make sure that Go is installed on your system. Clone the repository or download an archive for a specific version and run the following command in the terminal:

```bash
make build
```

This will create a binary file called `synctropy`; autocompletion files for `bash`, `zsh`, and `fish` in the directory `./autocompletions`; and markdown documentation files in the directory `./docs`. To install the binary and program files (autocompletion files and documentation):

```bash
sudo make install
```

And to uninstall:

```bash
sudo make uninstall
```

To remove files and directories created by this operation:

```bash
make clean
```

### Custom Destination Directory

By default, all files will be installed to their respectives subdirectories under the `/usr` directory in your root filesystem. However, if you want to set a custom destination directory for the installation, you can use the `DESTDIR` variable when running the `make install` command (also valid for `make uninstall`). For example, here's how you could build this program for Termux:

```bash
make install DESTDIR=/data/data/com.termux/files
```

### Initializing User Data Directory

Before you begin using `synctropy`, you should create the user data directory by executing the following command:

```bash
synctropy init
```

By default, this command will generate the necessary directory structure under `~/synctropy`. However, if you prefer a different location for your data directory, you can specify a custom path using the `-D` flag as shown below:

```bash
synctropy init -D <custom_data_dir>
```

## Usage

The general usage of Synctropy is as follows:

```
synctropy [command]
```

To get started, run the following command:

```
synctropy
```

This will display information about the program and provide instructions on how to proceed. To learn more about the available commands and options, you can use the `--help/-h` flag:

```
synctropy --help
```

This command will provide detailed information about the available commands, their usage, and the available options for each command. It is a useful reference when you need more information on how to use a specific command or what options are available for customization.

Feel free to explore the available commands and options using the `--help` flag to get a better understanding of the functionality provided by Synctropy.

### Documentation

To generate the program's documentation in markdown format, you can use the following command:

```shell
synctropy docs generate
```

The documentation will be generated in the `./docs` directory by default. You can specify a custom location by using:

```shell
synctropy docs generate -o <output_dir>
```

## Documentation

### Available Subcommands

- synctropy: The main command for the program.
  - version: Show the program's version.
  - init: Create user data directory
  - docs: Program documentation.
    - generate: Generate program documentation (markdown files).
  - utils: Utilities for hooks execution.
    - msg: Print a message in a specific color given a HEX code.
    - attention: Display attention message.
    - error: Display error message.
    - success: Display success message.
    - section: Display section title.
    - hr: Display a horizontal line.
    - confirm: Ask for confirmation.
    - sshagent-start: Starts an ssh-agent process and adds a private key.
    - sshagent-stop: Stops the ssh-agent process.
    - sshagent-getpid: Get the process ID of the ssh-agent.
    - sshagent-getsock: Get the socket path of the ssh-agent.
  - crates: Manage crates.
    - edit: Edit crates.
    - view: View crates.
    - create: Create crates.
    - rm: Remove crates.
    - ls: List crates.
    - hooks: Manage crate hooks.
      - run: Run crate hook(s).
      - ls: List crate hooks.
  - targets: Manage targets.
    - ls: List targets.
    - edit: Edit targets.
    - view: View targets.
    - sync: Sync targets.
    - enable: Enable targets.
    - disable: Disable targets.
    - create: Create targets.
    - rm: Remove targets.
    - hooks: Manage target hooks.
      - run: Run target hook(s).
      - ls: List target hooks.

### User Data Directory

The default user data directory is located at `~/synctropy`. Its tree structure is as follows:

- `templates`: This directory contains templates used for creating crates and targets. It provides a starting point with pre-configured setups for common synchronization scenarios. The templates are organized into subdirectories based on their type, such as `crates` and `targets`.

- `templates/crates`: This subdirectory within the `templates` directory contains templates specifically designed for creating crates. Each template may include a set of pre-defined hook scripts and configuration files to streamline the crate creation process.

- `templates/targets`: This subdirectory within the `templates` directory contains templates specifically designed for creating targets. Similar to the crate templates, each target template may include hook scripts and configuration files tailored for specific synchronization needs.

- `crates`: This directory holds the configurations and settings for all created crates. Each crate has its own subdirectory within the `crates` directory. The subdirectories are named after the respective crate and contain the associated configuration files, hooks, and any other necessary files.

- `crates/<crate>/targets`: Within each crate's subdirectory, there is a `targets` directory. This directory holds the configurations and hooks for all the targets associated with that particular crate. Each target has its own subdirectory within the `targets` directory, containing the target-specific configuration files, hooks, and any other necessary files.

#### Custom User Data Directory

It is possible to set a custom user data directory using the `-D` flag. To do so, you can run the `synctropy` program followed by the flag and the desired directory path. The command would look like this:

```
synctropy -D <custom_data_dir>
```

This flag can be used with any subcommand of the `synctropy` program. Whether you are creating crates, managing hooks, or performing other operations, you have the flexibility to specify a custom data directory that best suits your needs.

### Crates

`Crates` are the primary structural element in this program. A crate represents a collection of configurations and settings for syncing specific data (which are called `targets`).

Think of a crate as a container that holds the necessary information to establish a synchronization connection. Within a crate, you can define one or more targets, each corresponding to a specific configuration (e.g., Unison profiles) you want to synchronize.

#### Creating a Crate

To create a new crate, you can utilize the `crates create` command. This command will guide you through the process by prompting for a crate name and allowing you to choose a template. Templates provide pre-configured setups for common synchronization scenarios.

If there are no crate templates available or if you prefer to have a minimal setup with just the crate directory and a `targets` subdirectory, you can select the `scratch` template. This template will only create the necessary directories and will not include any additional configuration files or hooks.

After creating the crate, a new directory will be generated specifically for that crate. This directory will contain the relevant configuration files and hooks, based on the selected template.

#### Managing Crates

Once you have created crates, you can perform various operations on them. The following commands are available for managing crates:

- `crates edit`: Edit crates.
- `crates view`: View crates.
- `crates rm`: Remove crates.
- `crates ls`: List crates.
- `crates hooks`: Manage crate hooks.
 - `crates hooks run`: Run crate hook(s).
 - `crates hooks ls`: List crate hooks.

#### Hooks

Hooks are an essential component of crates, allowing you to define custom actions or scripts to execute. The following default hooks can be associated with a crate:

- `pre_transaction`: Runs before a synchronization transaction.
- `post_transaction`: Executes after a synchronization transaction is completed.
- `post_create`: Runs after a crate is created. Can be used to further configure the crate configuration beyond only creating its directory (which is done automatically by the program).
- `pre_rm`: Executes before removing a crate.
- `ls`: Displays custom information for the crate when running `crates ls`.
- `edit`: Script to open the crate configuration when running `crates edit`.
- `view`: Script to open the crate configuration when running `crates view`.

##### Environment Variables

When running crate hooks, the following environment variables are available for your use:

- **PROGRAM_NAME**: The name of the program (`synctropy`).
- **DEFAULT_SHELL**: The default shell used by the program.
- **SYNCTROPY_EXEC**: The executable path of `synctropy`.
- **SYNCTROPY_UTILS**: The executable path of `synctropy utils`, providing access to the utility commands.
- **USER_DATA_DIR**: The user data directory used by `synctropy`.
- **USER_CRATES_DIR**: The directory where user crates are stored.
- **USER_TEMPLATES_DIR**: The directory where user templates are stored.
- **USER_CRATES_TEMPLATES_DIR**: The directory where user crate templates are stored.
- **USER_TARGETS_TEMPLATES_DIR**: The directory where user target templates are stored.
- **CRATE_NAME**: The name of the current crate.
- **CRATE_DIR**: The directory path of the current crate.
- **CRATE_HOOKS_DIR**: The directory path of the hooks within the current crate.
- **CRATE_TARGETS_DIR**: The directory path of the targets within the current crate.
- **CRATE_TEMP_DIR**: The temporary directory path specific to the current crate.

These environment variables provide useful information and paths that can be utilized within your crate hooks to customize the behavior and perform specific actions based on the current context.

##### Custom entry command

A custom entry command for a specific hook can be defined by creating a file called `<hook_name>.entry` in the crate's hooks directory. For example, if you want to run the `post_create` hook with the `fish` shell, you could do so by creating a `post_create.entry` file with the following content:

```bash
/usr/bin/fish
```

##### Mananing Hooks

For managing hooks at the crate level, you can use the following subcommands:

- `crates hooks ls`: List all hooks available for a specific crate. This command displays the names of the hooks and their associated entry commands (if any).

- `targets hooks run`: Run one or more hooks for the specified crates. You can also run hooks in a specific sequence by running, for example:

```bash
crates hooks run --crate <crate_name> --hook <first_hook> --hook <second_hook>
```

##### Temporary Directory

When running hooks for crates, a temporary directory is created for each crate. This temporary directory serve as a workspace for performing actions or modifications during the hook execution process. It is called `.tmp` and created within the crate directory itself. This crate-specific temporary directory provides a separate workspace for any temporary files or data specific to the individual crate and can be used for any shared temporary files or data required by the hook scripts across multiple targets within the same crate. It allows the hooks to operate within a context that is isolated to the crate directory.

Once the hook execution is complete, the temporary directory and its contents are typically cleaned up automatically. When running hooks via `crates hooks run`, this behaviour can be modified using two options (non-exclusive):

- `nocreatetemp`: By using this option, temporary directories will not be created before running the hook(s). This can be useful if you prefer to handle temporary directory creation manually or if your hook scripts do not require a separate workspace.

- `noremovetemp`: With this option, temporary directories will not be automatically removed after running the hook(s). This allows you to inspect or access the temporary directories and their contents after the hook execution has completed. It can be beneficial for debugging purposes or if you need to access the temporary files generated during the hook execution.

### Targets

`Targets` are the secondary structural element and represent a specific configuration used to synchronize your files. Each target is associated with a crate and can have its own set of configurations and hooks.

#### Creating a Target

To create a new target, you can use the `targets create` command. This command will guide you through the process by allowing you to select a parent crate and provide a name for the target. Additionally, you can choose a template to pre-configure the target's setup according to your requirements.

If there are no target templates available or if you prefer a minimal setup for your target, you can select the `scratch` template during the target creation process. This template will create only the target's directory. It will not include any additional configuration files or hooks, providing you with a clean slate to customize according to your needs.

#### Managing Targets

Once you have created targets, you can perform various operations on them. The following commands are available for managing targets:

- `targets ls`: List targets.
- `targets edit`: Edit targets.
- `targets view`: View targets.
- `targets sync`: Sync targets.
- `targets enable`: Enable targets.
- `targets disable`: Disable targets.
- `targets rm`: Remove targets.
- `targets hooks`: Manage target hooks.
 - `targets hooks run`: Run target hook(s).
 - `targets hooks ls`: List target hooks.

#### Disabled Targets

When executing the `targets sync` command, `synctropy` will check if the target is disabled before initiating the synchronization process for each target. If any disabled targets are encountered, `synctropy` will skip them without generating an error. It will proceed to sync the remaining selected targets, if any are available. In essence, `synctropy` gracefully handles disabled targets during the synchronization process, allowing for the successful synchronization of the remaining enabled targets.

#### Hooks

Similar to crates, targets also support hooks that allow you to define custom actions or scripts to execute at specific stages of the synchronization process. The default hooks available for targets include:

- `pre_transaction`: Runs before the synchronization transaction.
- `post_transaction`: Executes after the synchronization transaction is completed.
- `post_create`: Runs after a target is created. Can be used to further configure the target configuration beyond only creating its directory (which is done automatically by the program).
- `pre_rm`: Executes before removing a target.
- `ls`: Displays custom information for the target when running `targets ls`.
- `edit`: Script to open the target configuration when running `targets edit`.
- `view`: Script to open the target configuration when running `targets view`.
- `sync`: Performs the actual synchronization transaction.

##### Environment Variables

When running target hooks, the following environment variables are available for your use:

- **PROGRAM_NAME**: The name of the program (`synctropy`).
- **DEFAULT_SHELL**: The default shell used by the program.
- **SYNCTROPY_EXEC**: The executable path of `synctropy`.
- **SYNCTROPY_UTILS**: The executable path of `synctropy utils`, providing access to the utility commands.
- **USER_DATA_DIR**: The user data directory used by `synctropy`.
- **USER_CRATES_DIR**: The directory where user crates are stored.
- **USER_TEMPLATES_DIR**: The directory where user templates are stored.
- **USER_CRATES_TEMPLATES_DIR**: The directory where user crate templates are stored.
- **USER_TARGETS_TEMPLATES_DIR**: The directory where user target templates are stored.
- **CRATE_NAME**: The name of the current crate.
- **CRATE_DIR**: The directory path of the current crate.
- **CRATE_HOOKS_DIR**: The directory path of the hooks within the current crate.
- **CRATE_TARGETS_DIR**: The directory path of the targets within the current crate.
- **CRATE_TEMP_DIR**: The temporary directory path specific to the current crate.
- **TARGET_NAME**: The name of the current target.
- **TARGET_DIR**: The directory path of the current target.
- **TARGET_HOOKS_DIR**: The directory path of the hooks within the current target.
- **TARGET_TEMP_DIR**: The temporary directory path specific to the current target.

These environment variables provide useful information and paths that can be utilized within your target hooks to customize the behavior and perform specific actions based on the current context.

##### Custom entry command

A custom entry command for a specific hook can be defined by creating a file called `<hook_name>.entry` in the target's hooks directory. For example, if you want to run the `post_create` hook with the `fish` shell, you could do so by creating a `post_create.entry` file with the following content:

```bash
/usr/bin/fish
```

##### Managing Hooks

When working with targets, you can use the following subcommands to manage their hooks:

- `targets hooks ls`: List all hooks available for the specified targets. This command displays the names of the hooks and their associated entry commands (if any).

- `targets hooks run`: Run one or more hooks for the specified targets. You can also run hooks in a specific sequence by running, for example:

```bash
targets hooks run --crate <crate_name> --target <target_name> --hook <first_hook> --hook <second_hook>
```

##### Temporary Directory

When running hooks for targets, two temporary directories are created. These temporary directories serve as a workspace for performing actions or modifications during the hook execution process.

- The first one is the `.tmp` directory within the crate directory, similar to the crate hooks. This temporary directory is used for any shared temporary files or data required by the hook scripts across multiple targets within the same crate.

- Additionally, a second temporary directory called `.tmp` is created within the target directory itself. This target-specific temporary directory provides a separate workspace for any temporary files or data specific to the individual target. It allows the hooks to operate within a context that is isolated to the target directory.

Once the hook execution is complete, the temporary directories and their contents are typically cleaned up automatically. When running hooks via `targets hooks run`, this behaviour can be modified using two options (non-exclusive):

- `nocreatetemp`: By using this option, temporary directories will not be created before running the hook(s). This can be useful if you prefer to handle temporary directory creation manually or if your hook scripts do not require a separate workspace.

- `noremovetemp`: With this option, temporary directories will not be automatically removed after running the hook(s). This allows you to inspect or access the temporary directories and their contents after the hook execution has completed. It can be beneficial for debugging purposes or if you need to access the temporary files generated during the hook execution.

### Synchronization Process

When configuring hooks for targets and crates, it's important to note that the `sync` hook (for `targets`) is the only required hook. All other pre/post hooks (for both `crates` and `targets`) are optional and can be customized based on your specific needs. Here's a breakdown of the hooks involved in the synchronization process, in order of execution:

1. **Pre-Transaction Hook (Crate)**: The `pre_transaction` hook for the crate is optional. It allows you to perform any necessary setup or checks that are applicable to all targets within the crate. This hook runs once at the beginning of the synchronization process for the entire crate.

2. **Pre-Transaction Hook (Target)**: The `pre_transaction` hook for each target is optional. It allows you to perform any target-specific setup or checks before the synchronization of that particular target. This hook runs before syncing each target.

3. **Sync Hook (Target)**: The `sync` hook is the core of the synchronization process for each target. It contains the necessary logic to synchronize your files. This hook is required for each target and must be defined to perform the actual synchronization.

4. **Post-Transaction Hook (Target)**: The `post_transaction` hook for each target is optional. It allows you to perform any cleanup or additional actions specific to that target after the synchronization is completed. This hook runs after syncing each target.

5. **Post-Transaction Hook (Crate)**: The `post_transaction` hook for the crate is optional. It allows you to perform any necessary cleanup or additional actions at the crate level after all selected targets have been synced. This hook runs once at the end of the synchronization process for the entire crate.

While the `sync` hook is required for each target, the other hooks provide flexibility to customize the synchronization process based on your specific requirements. You can choose to define and use the optional hooks as needed to perform additional actions or implement custom logic before and after syncing.

### Utilities

`synctropy` provides a set of utilities designed to be used within the hooks of crates and targets, allowing you to perform additional actions or execute custom logic during synchronization, though they can be used wherever and whenever you want. The main difference is that when running crate and target hooks, an environment variable called `$SYNCTROPY_UTILS` is automatically created, pointing to `synctropy utils`.

#### Available Utilities

- **msg**: Print a message in a specific color given a HEX code.
  ```
  synctropy utils msg '#FF5733' 'This is a colored message'
  ```

- **attention**: Display an attention message.
  ```
  synctropy utils attention 'This is an attention message'
  ```

- **error**: Display an error message.
  ```
  synctropy utils error 'This is an error message'
  ```

- **success**: Display a success message.
  ```
  synctropy utils success 'This is a success message'
  ```

- **section**: Display a section title.
  ```
  synctropy utils section 'This is a section title'
  ```

- **hr**: Display a horizontal line.
  ```
  synctropy utils hr '-' 0.45
  ```

- **confirm**: Ask for confirmation from the user.
  ```
  synctropy utils confirm 'Are you sure you want to continue?'
  ```

- **sshagent-start**: Start an ssh-agent process and add a private key.
  ```
  synctropy utils sshagent-start '/path/to/my/key' $TARGET_TEMP_DIR
  ```

- **sshagent-stop**: Stop the ssh-agent process.
  ```
  synctropy utils sshagent-stop $TARGET_TEMP_DIR
  ```

- **sshagent-getpid**: Get the process ID of the ssh-agent.
  ```
  synctropy utils sshagent-getpid $TARGET_TEMP_DIR
  ```

- **sshagent-getsock**: Get the socket path of the ssh-agent.
  ```
  synctropy utils sshagent-getsock $TARGET_TEMP_DIR
  ```

#### Using Utilities in Hooks

To use any of the utilities within a crate or target hook, you can access them using the `$SYNCTROPY_UTILS` environment variable, which points to the command `synctropy utils`. For example, to display an attention message within a crate hook:

```bash
$SYNCTROPY_UTILS attention 'This is an attention message'
```

This will print the specified attention message to the console, attracting the user's attention during the synchronization process.

### Templates

Templates play a crucial role in customizing the creation of crates and targets. `synctropy` provides a template-based approach to create crates and targets, allowing you to quickly set up and configure your project structure. When creating a crate or target, `synctropy` automatically generates the corresponding crate or target directory and copies the selected template structure to it.
To simplify template usage, you can find a collection of example crate and target templates that I personally use in the `./templates` directory within the source code. These templates serve as starting points and can be customized to suit your specific project requirements. Additionally, you have the flexibility to create your own templates and store them in the following directories within the user data directory:

- `templates/crates`: This directory is dedicated to crate templates.
- `templates/targets`: This directory is dedicated to target templates.

By placing your custom templates in these directories, they become readily available for selection during the crate and target creation process. You can leverage these templates to expedite the setup of your projects and tailor them to your specific needs.

## License

Synctropy is licensed under the GPL-3.0 license.
