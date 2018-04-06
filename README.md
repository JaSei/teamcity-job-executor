# teamcity job executor

[![Release](https://img.shields.io/github/release/avast/teamcity-job-executor.svg?style=flat-square)](https://github.com/avast/teamcity-job-executor/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Linux build status](https://img.shields.io/travis/avast/teamcity-job-executor.svg?style=flat-square)](https://travis-ci.org/avast/teamcity-job-executor)

simple command line utility which runs job in teamcity and waits (without --nowait flag)
for job to finish

success status of build = exit code 0

other status = exit code 1

non-interactive run example:
```
teamcity-job-executor -H teamcity.server -u my_username -p my_password backends_Modules_App_Test
```

help:
```
usage: teamcity-job-executor --hostname=HOSTNAME --username=USERNAME [<flags>] [<configId>]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -H, --hostname=HOSTNAME  teamcity hostname
  -u, --username=USERNAME  teamcity username
  -p, --password=PASSWORD  teamcity password
   -j, --job_param=JOB_PARAM ...  teamcity job parameters in key=value format
      --sleep=5s           sleep duration of pooling teamcity
      --nowait             Does not wait for queued job to finish
      --version            Show application version.

Args:
  [<configId>]  id of build configuration which you can run
```
