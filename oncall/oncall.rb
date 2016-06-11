module Oncall
  class Bot < SlackRubyBot::Bot

    def self.slack
      @@slack ||= Slack::Web::Client.new
      @@slack
    end


    def self.route
      
    end

  end
end
