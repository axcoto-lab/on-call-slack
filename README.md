# Oncall

Automatically runs fix via slack interaction

Currently we trigger manually. however, we can configure bot to listen
to message like a pagerduty message

# How it works

It invoke command via SSH, all are shell script

# How to define a fix

Create a file in `fixer` folder with this structure

```ruby
module Oncall
  module Fixer
    class Name
      include Oncall::Router
      include Oncall::Executor

      # Define how do we SSH into jump host
      ssh do
        {host: host, username: username_to_ssh; ssh_key: path to key}
      end

      # Define mapping command
      route do
        command 'fix_something' do |client, data, match|
          say "I'm processing #{match[:some_argument]}"
          result = run("this command is run over SSH connection above")
          output.append result
          # Or we can do some kind of processing
          post_result = crazy_processing result
          # Then we can continue to run another command 
          result = run("another command based on post_result")
          output.append result
        end
      end
    end

  end
end
```
