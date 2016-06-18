module Oncall
  class Bot < SlackRubyBot::Bot
    # Get an instance of slack to use through out of the app
    def self.slack
      @@slack ||= Slack::Web::Client.new
      @@slack
    end
  end


  # Define router
  module Router
    class Builder

    end

    def self.instance
      @@instance = Builder.new
    end

    def self.route
      
    end

    def self.included(base)
      
    end
  end

  # Define how we run command
  module Executor
    def self.ssh
      @@ssh = yield
    end

    # Run a command remove via jumphost
    def self.run_remote

    end

    # Run a command locally
    def self.local

    end
  end
end
