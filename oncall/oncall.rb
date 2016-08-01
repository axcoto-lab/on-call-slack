require_relative '../lib/wit'

module Oncall
  class Bot < SlackRubyBot::Bot
    include Oncall::Wit

    # Get an instance of slack to use through out of the app
    def self.slack
      @@slack ||= Slack::Web::Client.new
      @@slack
    end

    match /^.*vinhbot\s+(.*)/ do |client, data, match|
      puts "Will parse #{match[1]}"
      message = wit.parse match[1]
      puts message.to_a.to_s
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
