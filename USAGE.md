## hashstack admin agent-delete

Deletes an agent by id.

### Synopsis



Delete an agent by id. This does not prevent the agent from
communicating with the server. This is only useful for when an
agent is brought offline and will never return.
    

```
hashstack admin agent-delete <agent_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin impersonate

Create an authentication token for another user.

### Synopsis



Create an authentication token for another user by username. This
can be useful when troubleshooting as another user.
    

```
hashstack admin impersonate <username>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin job-delete

Deletes a job by project_id and job_id.

### Synopsis


Deletes a job by project_id and job_id.

```
hashstack admin job-delete <project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin job-pause

Pauses a job by project_id and job_id.

### Synopsis


Pauses a job by project_id and job_id.

```
hashstack admin job-pause <project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin jobs

Displays a list of all active jobs.

### Synopsis


Displays a list of all active jobs.

```
hashstack admin jobs
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin job-start

Starts a job by project_id and job_id.

### Synopsis


Starts a job by project_id and job_id.

```
hashstack admin job-start <project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack admin

Executes admin subcommands (-h or --help for more info).

### Synopsis



Executes a subcommand as an administrator (-h or --help for a list of subcommands.
    

```
hashstack admin
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack admin agent-delete](hashstack_admin_agent-delete.md)	 - Deletes an agent by id.
* [hashstack admin impersonate](hashstack_admin_impersonate.md)	 - Create an authentication token for another user.
* [hashstack admin job-delete](hashstack_admin_job-delete.md)	 - Deletes a job by project_id and job_id.
* [hashstack admin job-pause](hashstack_admin_job-pause.md)	 - Pauses a job by project_id and job_id.
* [hashstack admin job-start](hashstack_admin_job-start.md)	 - Starts a job by project_id and job_id.
* [hashstack admin jobs](hashstack_admin_jobs.md)	 - Displays a list of all active jobs.
* [hashstack admin projects](hashstack_admin_projects.md)	 - Displays a list of all projects.


## hashstack admin projects

Displays a list of all projects.

### Synopsis


Displays a list of all projects.

```
hashstack admin projects
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).


## hashstack agents

Display a list of agents connected to the cluster.

### Synopsis



Display a list of agents connected to the cluster. If an id is provided, only information
on that agent will be displayed.

```
hashstack agents [id]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack hcstats add

Upload the provided file to the server.

### Synopsis


Upload the provided file to the server.

```
hashstack hcstats add <file>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack hcstats](hashstack_hcstats.md)	 - Display a list of all hcstat files available on the server (-h or --help for subcommands).


## hashstack hcstats delete

Delete a file by name from the server.

### Synopsis


Delete a file by name from the server.

```
hashstack hcstats delete <file_name>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack hcstats](hashstack_hcstats.md)	 - Display a list of all hcstat files available on the server (-h or --help for subcommands).


## hashstack hcstats

Display a list of all hcstat files available on the server (-h or --help for subcommands).

### Synopsis



Displays a list of hcstat files that are stored on the remote server. If file_name is provided, details will be displayed for that specific
hcstat file.

hcstat files can be used in jobs. Additional subcommands are available to add and delete files.


```
hashstack hcstats [file_name]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack hcstats add](hashstack_hcstats_add.md)	 - Upload the provided file to the server.
* [hashstack hcstats delete](hashstack_hcstats_delete.md)	 - Delete a file by name from the server.


## hashstack jobs add

Add a job for the provided project and list.

### Synopsis


Add a job for the provided project and list.

```
hashstack jobs add <project_name|project_id> <list_name|list_id> <name> <wordlist|mask>
```

### Options

```
  -a, --attack-mode int           Attack-mode, see references below
  -1, --custom-charset1 string    User-defined charset ?1
  -2, --custom-charset2 string    User-defined charset ?2
  -3, --custom-charset3 string    User-defined charset ?3
  -4, --custom-charset4 string    User-defined charset ?4
      --hex-charset               Assume charset is given in hex
      --markov-hcstat string      Specify hcstat file to use
  -t, --markov-threshold int      Threshold X when to stop accepting new markov-chains
      --max-devices int           Maximum devices across the entire cluster to use, 0 is unlimited
      --opencl-vector-width int   Manual override OpenCL vector-width to X
      --priority int              The priority for this job 1-100 (default 1)
  -j, --rule-left string          Single rule applied to each word from left wordlist
  -k, --rule-right string         Single rule applied to each word from left wordlist
  -r, --rules-file string         Rule file to be applied to each word from wordlists
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).


## hashstack jobs delete

Deletes a job by project_name|project_id and job_id.

### Synopsis


Deletes a job by project_name|project_id and job_id.

```
hashstack jobs delete <project_name|project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).


## hashstack jobs

Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).

### Synopsis


Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).

```
hashstack jobs <project_name|project_id> [job_id]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack jobs add](hashstack_jobs_add.md)	 - Add a job for the provided project and list.
* [hashstack jobs delete](hashstack_jobs_delete.md)	 - Deletes a job by project_name|project_id and job_id.
* [hashstack jobs pause](hashstack_jobs_pause.md)	 - Pauses a job by project_name|project_id and job_id.
* [hashstack jobs start](hashstack_jobs_start.md)	 - Starts a job by project_name|project_id and job_id.
* [hashstack jobs update](hashstack_jobs_update.md)	 - Updates a job by project_name|project_id and job_id.


## hashstack jobs pause

Pauses a job by project_name|project_id and job_id.

### Synopsis


Pauses a job by project_name|project_id and job_id.

```
hashstack jobs pause <project_name|project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).


## hashstack jobs start

Starts a job by project_name|project_id and job_id.

### Synopsis


Starts a job by project_name|project_id and job_id.

```
hashstack jobs start <project_name|project_id> <job_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).


## hashstack jobs update

Updates a job by project_name|project_id and job_id.

### Synopsis



Updates a job by project_name|project_id and job_id. Can be used to update
priority and/or max-devices.
	

```
hashstack jobs update <project_name|project_id> <job_id>.
```

### Options

```
      --max-devices int   Maximum devices across the entire cluster to use, 0 is unlimited
      --priority int      The priority for this job 1-100 (default 1)
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).


## hashstack lists add

Add a new list to a project.

### Synopsis



Add a new file containing one or more hashes to a project by project_name or project_id. Modes can be viewed
using the "modes" subcommand. The file name must be unique across projects.


```
hashstack lists add <project_name|project_id> <mode> <file>
```

### Options

```
      --hex-salt   Assume is given in hex
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack lists](hashstack_lists.md)	 - Displays a list of all lists associated with the provided project (-h or --help for subcommands).


## hashstack lists delete

Delete a list from a project.

### Synopsis



Delete a list from a project by project_name or project_id. Deleting a list also deletes
the associated plains.


```
hashstack lists delete <project_name|project_id> <list_name|list_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack lists](hashstack_lists.md)	 - Displays a list of all lists associated with the provided project (-h or --help for subcommands).


## hashstack lists hashes

Download hashes for a list.

### Synopsis


Download hashes for a list.

```
hashstack lists hashes <project_name|project_id> <list_name|list_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack lists](hashstack_lists.md)	 - Displays a list of all lists associated with the provided project (-h or --help for subcommands).


## hashstack lists

Displays a list of all lists associated with the provided project (-h or --help for subcommands).

### Synopsis



Displays a list of all lists associated with the provided project. If list_name|list_id is provided, details will be displayed for
that specific list. Additional subcommands are available.


```
hashstack lists <project_name|project_id> [list_name|list_id]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack lists add](hashstack_lists_add.md)	 - Add a new list to a project.
* [hashstack lists delete](hashstack_lists_delete.md)	 - Delete a list from a project.
* [hashstack lists hashes](hashstack_lists_hashes.md)	 - Download hashes for a list.
* [hashstack lists plains](hashstack_lists_plains.md)	 - Download plains for a list.


## hashstack lists plains

Download plains for a list.

### Synopsis


Download plains for a list.

```
hashstack lists plains <project_name|project_id> <list_name|list_id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack lists](hashstack_lists.md)	 - Displays a list of all lists associated with the provided project (-h or --help for subcommands).


## hashstack login

Login and cache a session token for the remote server.

### Synopsis



This command will prompt for your password and send it along with your username to the server at server_url.
The token returned along with the server_url will be saved in your home directory for all additional requests.
    

```
hashstack login <server_url> <username>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack

Execute commands against a Hashstack server. Try -h or --help for more information.

### Synopsis


Execute commands against a Hashstack server. Try -h or --help for more information.

```
hashstack
```

### Options

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack admin](hashstack_admin.md)	 - Executes admin subcommands (-h or --help for more info).
* [hashstack agents](hashstack_agents.md)	 - Display a list of agents connected to the cluster.
* [hashstack hcstats](hashstack_hcstats.md)	 - Display a list of all hcstat files available on the server (-h or --help for subcommands).
* [hashstack jobs](hashstack_jobs.md)	 - Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands).
* [hashstack lists](hashstack_lists.md)	 - Displays a list of all lists associated with the provided project (-h or --help for subcommands).
* [hashstack login](hashstack_login.md)	 - Login and cache a session token for the remote server.
* [hashstack modes](hashstack_modes.md)	 - Prints a list of supported hash modes.
* [hashstack projects](hashstack_projects.md)	 - Display a list of all projects (-h or --help for subcommands).
* [hashstack rules](hashstack_rules.md)	 - Display a list of all rule files available on the server (-h or --help for subcommands).
* [hashstack status](hashstack_status.md)	 - Displays information about the Hashstack cluster.
* [hashstack users](hashstack_users.md)	 - Displays a list of hashstack users.
* [hashstack version](hashstack_version.md)	 - Print client and server version and exit.
* [hashstack wordlists](hashstack_wordlists.md)	 - Display a list of all wordlists available on the server (-h or --help for subcommands).


## hashstack modes

Prints a list of supported hash modes.

### Synopsis


Prints a list of supported hash modes.

```
hashstack modes
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack projects add-contributor

Adds a user to the project by username.

### Synopsis


Adds a user to the project by username.

```
hashstack projects add-contributor  <name|id> <username>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack projects](hashstack_projects.md)	 - Display a list of all projects (-h or --help for subcommands).


## hashstack projects add

Add a new project.

### Synopsis



Add a new project using the provided name and description. The name must be unique across all projects.
The name should not include special characters or spaces. The description is required.
	

```
hashstack projects add <project_name> <description>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack projects](hashstack_projects.md)	 - Display a list of all projects (-h or --help for subcommands).


## hashstack projects delete

Delete a project by name or id.

### Synopsis



Delete a project by name or id. Deleting a project will delete any associated lists.


```
hashstack projects delete <name|id>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack projects](hashstack_projects.md)	 - Display a list of all projects (-h or --help for subcommands).


## hashstack projects

Display a list of all projects (-h or --help for subcommands).

### Synopsis



Displays a list of your projects. If name or id is provided, details will be displayed for that specific project.
Additional subcommands are available.

Projects are used to organize one or more lists. An exmaple of a project may be AcmeForensicIvestigation2016.
Once a project is created, you will upload your lists using the project's name. You can share projects with many users.


```
hashstack projects [name|id]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack projects add](hashstack_projects_add.md)	 - Add a new project.
* [hashstack projects add-contributor](hashstack_projects_add-contributor.md)	 - Adds a user to the project by username.
* [hashstack projects delete](hashstack_projects_delete.md)	 - Delete a project by name or id.
* [hashstack projects remove-contributor](hashstack_projects_remove-contributor.md)	 - Removes a user to the project by username


## hashstack projects remove-contributor

Removes a user to the project by username

### Synopsis


Removes a user to the project by username.

```
hashstack projects remove-contributor  <name|id> <username>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack projects](hashstack_projects.md)	 - Display a list of all projects (-h or --help for subcommands).


## hashstack rules add

Upload the provided file to the server.

### Synopsis


Upload the provided file to the server.

```
hashstack rules add <file>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack rules](hashstack_rules.md)	 - Display a list of all rule files available on the server (-h or --help for subcommands).


## hashstack rules delete

Delete a file by name from the server.

### Synopsis


Delete a file by name from the server.

```
hashstack rules delete <file_name>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack rules](hashstack_rules.md)	 - Display a list of all rule files available on the server (-h or --help for subcommands).


## hashstack rules

Display a list of all rule files available on the server (-h or --help for subcommands).

### Synopsis



Displays a list of rule files that are stored on the remote server. If file_name is provided, details will be displayed for that specific
rule file.

Rule files can be used in jobs. Additional subcommands are available to add and delete files.


```
hashstack rules [file_name]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack rules add](hashstack_rules_add.md)	 - Upload the provided file to the server.
* [hashstack rules delete](hashstack_rules_delete.md)	 - Delete a file by name from the server.


## hashstack status

Displays information about the Hashstack cluster.

### Synopsis


Displays information about the Hashstack cluster.

```
hashstack status
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack users

Displays a list of hashstack users.

### Synopsis



You may list users here. Interactions with hashstack for other operations
should be conducted directly on the server by an administrator.

Users can be added or removed from a project using the projects command.
    

```
hashstack users
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack version

Print client and server version and exit.

### Synopsis


Print client and server version and exit.

```
hashstack version
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.


## hashstack wordlists add

Upload the provided file to the server.

### Synopsis


Upload the provided file to the server.

```
hashstack wordlists add <file>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack wordlists](hashstack_wordlists.md)	 - Display a list of all wordlists available on the server (-h or --help for subcommands).


## hashstack wordlists delete

Delete a file by name from the server.

### Synopsis


Delete a file by name from the server.

```
hashstack wordlists delete <file_name>
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack wordlists](hashstack_wordlists.md)	 - Display a list of all wordlists available on the server (-h or --help for subcommands).


## hashstack wordlists

Display a list of all wordlists available on the server (-h or --help for subcommands).

### Synopsis



Displays a list of wordlists that are stored on the remote server. If file_name is provided, details will be displayed for that specific
wordlist.

Wordlists can be used in jobs. Additional subcommands are available to add and delete files.


```
hashstack wordlists [file_name]
```

### Options inherited from parent commands

```
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```

### SEE ALSO
* [hashstack](hashstack.md)	 - Execute commands against a Hashstack server. Try -h or --help for more information.
* [hashstack wordlists add](hashstack_wordlists_add.md)	 - Upload the provided file to the server.
* [hashstack wordlists delete](hashstack_wordlists_delete.md)	 - Delete a file by name from the server.


