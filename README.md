# on-call
Auto detect command to run when doing on call

# Why

I recently started doing on-call, and it sucks. I don't have many
experience with some software components that I have to monitor. Example
I never used Haddop/Hive personally. So I don't have a good
understanding of it and failed to really troubleshoot it other than
blindly restart a broken system.

So I write this blindness tool, it will fetch information from
pagerduty, and compare with a [hjson](hjson.org) file which define what
to do.

I simply run this tool over it.

