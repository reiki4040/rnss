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

homebrew (NOT YET)
```
brew install reiki4040/tap/rnss
```

### other required

installed `aws cli` and `session-manage-plugin`. 
- [aws cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- [session-manager-plugin](https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

## usage

```
rnss [options] <filter phrase>
```

show instances and select then start-session to the instance.


