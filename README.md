# hashstack-cli
A cross-platform CLI interface to Hashstack


## Usage
#### General
```
Execute commands against a Hashstack server. Try -h or --help for more information.

Usage:
  hashstack-cli [flags]
  hashstack-cli [command]

Available Commands:
  agents      Display a list of agents connected to the cluster
  hcstats     Display a list of all hcstat files available on the server (-h or --help for subcommands)
  jobs        Display a list of jobs for a project or attach to a job by id (-h or --help for subcommands)
  lists       Displays a list of all lists associated with the provided project (-h or --help for subcommands)
  login       Login and cache a session token for the remote server
  modes       Prints a list of supported hash modes
  projects    Display a list of all projects (-h or --help for subcommands)
  rules       Display a list of all rule files available on the server (-h or --help for subcommands)
  status      Displays information about the Hashstack cluster
  users       Displays a list of hashstack users
  version     Print the version and exit
  wordlists   Display a list of all wordlists available on the server (-h or --help for subcommands)

Flags:
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation

Use "hashstack-cli [command] --help" for more information about a command.
```

#### Add Job
```
.\hashstack-cli.exe jobs add --help
Add a job for the provided project and list

Usage:
  hashstack-cli jobs add <project_name|project_id> <list_name|list_id> <name> <dictionary|mask> [flags]

Flags:
  -a, --attack-mode int           Attack-mode, see references below
  -1, --custom-charset1 string    User-defined charset ?1
  -2, --custom-charset2 string    User-defined charset ?2
  -3, --custom-charset3 string    User-defined charset ?3
  -4, --custom-charset4 string    User-defined charset ?4
      --hex-charset               Assume charset is given in hex
      --markov-hcstat string      Specify hcstat file to use
  -t, --markov-threshold int      Threshold X when to stop accepting new markov-chains
      --max-devices int           Maximum devices across the entire cluster to use, 0 is unlimited
      --opencl-vector-width int   Manual override OpenCL  vector-width to X
      --priority int              The priority for this job 1-100 (default 1)
  -j, --rule-left string          Single rule applied to each word from left wordlist
  -k, --rule-right string         Single rule applied to each word from left wordlist
  -r, --rules-file string         Rule file to be applied to each word from wordlists

Global Flags:
      --config string   config file (default: $HOME/.hashstack/config)
      --debug           enable debug output
      --insecure        skip TLS certificate validation
```      


## Start to Finish
#### Login
```
 .\hashstack-cli.exe login http://192.168.7.98:8000 admin
Password: *****
Authentication credentials cached in C:\Users\tom\.hashstack\config
```

#### Add a project
```
 .\hashstack-cli.exe projects add x "this is x"
ID...............: 4
Name.............: x
Jobs.............: 0/0 Active 0/0 Complete
Lists............: 0
Owner............: admin
Last Updated.....: now
```

#### Add a list to that project
```
.\hashstack-cli.exe lists add x 0 C:\Users\tom\Downloads\hashcat-3.30\example0.hash
ID..............: 2
Name............: example0.hash
Hash Mode.......: 0
Recovered........: 0/6494 (0%)
```

#### Start a job for that list
```
.\hashstack-cli.exe jobs add x example0.hash brute8 -a 3 "?a?a?a?a?a?a"
Job.............: brute8
ID..............: 5
Status..........: Running
Hash.Type.......: 0 (MD5)
Hash.Target.....: example0.hash
Max Devices.....: 0
Priority........: 1
Time.Created....: 5 seconds ago
Time.Started....: a long while ago
Recovered.......: 0/6494 (0%)
```

#### Get Plains
```
.\hashstack-cli.exe lists plains x example0.hash
e11c594e6a2f4eb499cceadfca988595:13LEXON
```

## Cluster Stats
```
.\hashstack-cli.exe status
Jobs...........: 0 Active, 0 Paused
Node.Count.....: 1
GPU.Count......: 1
CPU.Count......: 0
Load.GPU.......: 0/100 (0%)
Load.CPU.......: 0/0 (0%)
Temp.GPU.......: 35C - 35C
Temp.CPU.......: 0C - 0C
```

## Agents
```
.\hashstack-cli.exe agents
ID..............: 1
Host............: ht
IP.Address......: 192.168.224.13
Uptime..........: 2 days 2 hours 57 minutes 6 seconds
Last.Seen.......: 3 days ago
Memory..........: 2.4 GB/17 GB (15%)
Dev.#1..........: GeForce GTX 750 Ti, 1032 Mhz, 0% load, 35C, 33% Fan
```

## Supported Modes
```
.\hashstack-cli.exe modes
#       Name
0       MD5
2500    WPA/WPA2
9700    Office 20
```

## Rules/Wordlists
```
.\hashstack-cli.exe rules
Name.............: Incisive-leetspeak.rule
Rules............: 15487
Size.............: 325 kB
Last Modified....: 4 days ago

Name.............: InsidePro-HashManager.rule
Rules............: 6746
Size.............: 42 kB
Last Modified....: 4 days ago
```
