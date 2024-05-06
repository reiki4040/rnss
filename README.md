rnss
===
connect EC2 using session manager helper.

rnss is instance selection helper for ssm start session.
you can show EC2 instances and select in CUI then start ssm sesion.
rnss is simple wrapper that call below.

```
aws ssm start-session --target <you selected instance-id>
```

## install

homebrew
```
brew install reiki4040/tap/rnss
```

*if you get like below message, then run `brew update`*

```
Warning: No available formula or cask with the name "reiki4040/tap/rnss". Did you mean reiki4040/tap/rnssh, reiki4040/tap/rnsd or reiki4040/tap/rnzoo?
```

this is typo suggest from local stored tap info that not included rnss.

```
brew update
```

### other required

installed `aws cli` and `session-manager-plugin`. 
- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

## usage

show instances and select then start-session to the instance.

```
rnss [-region AWS_REGION] [-f] [filter phrase...]
```

set AWS region using `AWS_REGION` enviroment variable or `-region`

### basic usage

show ec2 list (with cache if exists) and connect
```
rnss -region ap-northeast-1
```

without cache (refresh cache, call ec2 describe instances and store local file). set `-f` if you changed EC2 state ex. launch instance, stop instance.
```
rnss -region ap-northeast-1 -f
```

command args is filter phrase. filter `web server` in first list view by below.
```
rnss -region ap-northeast-1 web server
```

## run example

local run example that using `-stdin` and `-show-command` and initial filter `web server`.

```
cat example_ec2list.txt | rnss -stdin -show-command web server
```

- `-stdin`: ec2 list from stdin. this example example_ec2list.txt from pipe.
- `-show-command`: only show command, NOT run start session. it likes dry-run.

show instance list and select any(press Enter) show `aws ssm start-session --target <your selected InstanceId>`


