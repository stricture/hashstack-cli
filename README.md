# hashstack-cli
A cross-platform CLI interface to Hashstack

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