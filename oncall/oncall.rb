require_relative '../lib/wit'

module Oncall
  class Bot < SlackRubyBot::Bot
    extend Oncall::Wit

    # Get an instance of slack to use through out of the app
    def self.slack
      @@slack ||= Slack::Web::Client.new
      @@slack
    end

    match /^.*vinhbot\s+(.*)/ do |client, data, match|
      message = wit.parse match[1]

      if message["aws_resource"]
        case message["aws_attribute"]
          when "ip"
            puts "Will find #{message["aws_attribute"]} of#{message["aws_resource"]}"
            default = `aws ec2 describe-instances --filters Name=tag-key,Values="*#{message["aws_resource"]}*" --query 'Reservations[].Instances[*].[NetworkInterfaces[].PrivateIpAddresses[].PrivateIpAddress, Tags[1].Value]' --profile cc --output text`
            eu = `aws --region eu-central-1 ec2 describe-instances --filters Name=tag-key,Values="*#{message["aws_resource"]}*" --query 'Reservations[].Instances[*].[NetworkInterfaces[].PrivateIpAddresses[].PrivateIpAddress, Tags[1].Value]' --profile cc --output text`
            slack.chat_postMessage(text: "```" + [default, eu].join("\n") + "```", channel: data.channel)
        end
      end

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
